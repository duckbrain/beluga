package beluga

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/duckbrain/beluga/internal/compose"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Runner struct {
	Environment Environment
}

func New() Runner {
	return Runner{
		Environment: Env(),
	}
}

func (r Runner) Exec(c *exec.Cmd) error {
	if c.Stdin == nil {
		c.Stdin = os.Stdin
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env, err = r.Environment.Format(GoEnvFormat, true)
	if err != nil {
		return errors.Wrap(err, "generate environment list")
	}
	return c.Run()
}

func (r Runner) ComposeFile(ctx context.Context) (string, error) {
	cmd := exec.Command("docker-compose", "config")
	cmd.Env, _ = env.Format(GoEnvFormat, true)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	file := compose.File{}
	err = yaml.Unmarshal(out, &file)
	if err != nil {
		return "", errors.Wrap(err, "unmarshal compose output")
	}

	mod := r.serviceModifier()

	for name, service := range file.Services {
		port, _ := strconv.ParseUInt(service.Labels.Get("us.duckfam.beluga.port"), 10, 16)

		info := serviceInfo{
			Port: uint16(port),
		}
		mod(name, &service, info)
		file.Services[name] = service
	}

	data, err := yaml.Marshal(file)
	return string(data), errors.Wrap("marshal compose file", err)
}

type serviceInfo struct {
	Port uint16
}
type serviceModifier func(name string, service *compose.Service, info serviceInfo)

func (r Runner) serviceModifier() serviceModifier {
	host := r.Environment.Domain()

	return func(name string, service *compose.Service, info serviceInfo) {
		if info.Port == 0 {
			return
		}

		service.Environment.Set("VIRTUAL_HOST", host)
		service.Environment.Set("VIRTUAL_PORT", fmt.Sprintf("%v", info.Port))
	}
}
