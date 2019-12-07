package beluga

import (
	"context"
	"io/ioutil"
	"net/url"
	"os/exec"
)

type DeployMode int

const (
	ComposeMode DeployMode = 2
	SwarmMode   DeployMode = 1
)

type DeployOpts struct {
	ComposeFile string
	StackName   string
	Mode        DeployMode
}

func (opts DeployOpts) writeComposeFile() (filename string, err error) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	_, err = file.WriteString(opts.ComposeFile)
	if err != nil {
		return "", err
	}
	return file.Name(), file.Close()
}

type Deployer interface {
	Deploy(opts DeployOpts) error
	Teardown(opts DeployOpts) error
}

func (env Environment) ComposeFile(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "docker-compose", "config")
	cmd.Env, _ = env.Format(GoEnvFormat, true)
	out, err := cmd.Output()
	return string(out), err
}

func (env Environment) Deployer() Deployer {
	host := env[varDeployDockerHost]
	u, _ := url.Parse(host)
	switch u.Scheme {
	case "portainer":
		return &Portainer{DSN: u}
	case "ssh":
		panic("SSH not implemented")
	default:
		return Docker(host)
	}
}

func (env Environment) DeployOpts() DeployOpts {
	return DeployOpts{
		ComposeFile: env[varDockerComposeFile],
		StackName:   env[varStackName],
		Mode:        ComposeMode,
	}
}
