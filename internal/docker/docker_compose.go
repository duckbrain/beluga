package docker

import (
	"github.com/duckbrain/beluga/internal/lib"
)

type Deployer interface {
	Deploy(composeFile string) error
	Teardown(composeFile string) error
}
type Compose lib.Environment

func (c Compose) run(args ...string) error {
	e := lib.Environment(c)
	a := []string{}
	if v := e.DockerComposeFile(); v != "" {
		a = append(a, "--file", v)
	}
	if v := e.DeployDockerHost(); v != "" {
		a = append(a, "--host", v)
	}
	a = append(a, args...)
	return run("docker-compose", a...)
}

func (c Compose) Deploy() error {
	return c.run("up", "--detach", "--no-build")
}

func (c Compose) Teardown() error {
	return c.run("down")
}
