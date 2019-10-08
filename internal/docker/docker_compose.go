package docker

import (
	"os/exec"

	"github.com/duckbrain/beluga/internal/lib"
)

type Deployer interface {
	Deploy() error
	Teardown() error
}
type Compose struct {
	env lib.Environment
	runner
}

func (c Compose) run(args ...string) error {
	e := c.env
	a := []string{}
	if v := e.DockerComposeFile(); v != "" {
		a = append(a, "--file", v)
	}
	if v := e.DeployDockerHost(); v != "" {
		a = append(a, "--host", v)
	}
	a = append(a, args...)
	return c.runner.run(exec.Command("docker-compose", a...))
}

func (c Compose) Deploy() error {
	return c.run("up", "--detach", "--no-build")
}

func (c Compose) Teardown() error {
	return c.run("down")
}
