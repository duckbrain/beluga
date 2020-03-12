package beluga

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"text/template"

	"github.com/duckbrain/beluga/internal/compose"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Runner struct {
	Env Environment
}

func New() Runner {
	return Runner{
		Env: Env(),
	}
}

// Exec runs a command in the Beluga context with stderr and stdout bound to the parent process
func (r Runner) Exec(c *exec.Cmd) error {
	var err error
	if c.Stdin == nil {
		c.Stdin = os.Stdin
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env, err = r.Env.Format(GoEnvFormat, true)
	if err != nil {
		return errors.Wrap(err, "generate environment list")
	}
	return c.Run()
}

func (r Runner) ComposeFile(ctx context.Context) (string, error) {
	cmd := exec.Command("docker-compose", "config")
	cmd.Env, _ = r.Env.Format(GoEnvFormat, true)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	templateS := r.Env.DockerComposeTemplate()
	if len(templateS) == 0 {
		return string(out), nil
	}
	templateBase := compose.File{}
	err = yaml.Unmarshal([]byte(templateS), &templateBase)
	if err != nil {
		return "", errors.Wrap(err, "parse compose template yaml")
	}

	t, err := template.New("").Parse(templateS)
	if err != nil {
		return "", errors.Wrap(err, "parse compose template")
	}

	file := compose.File{}
	err = yaml.Unmarshal(out, &file)
	if err != nil {
		return "", errors.Wrap(err, "unmarshal compose output")
	}

	file.Fields.Merge(templateBase.Fields)

	for name, service := range file.Services {
		port, _ := strconv.ParseUint(service.Labels.Get("us.duckfam.beluga.port"), 10, 16)

		info := serviceInfo{
			Port: uint16(port),
		}

		s := new(bytes.Buffer)
		err := t.Execute(s, info)
		if err != nil {
			return "", errors.Wrap(err, "execute compose template")
		}

		file.Services[name] = service
	}

	data, err := yaml.Marshal(file)
	return string(data), errors.Wrap(err, "marshal compose file")
}

type serviceInfo struct {
	Port uint16
}

func (r Runner) serviceModifier() serviceModifier {
	host := r.Env.Domain()

	return func(name string, service *compose.Service, info serviceInfo) {
		if info.Port == 0 {
			return
		}

		service.Environment.Set("VIRTUAL_HOST", host)
		service.Environment.Set("VIRTUAL_PORT", fmt.Sprintf("%v", info.Port))
	}
}
