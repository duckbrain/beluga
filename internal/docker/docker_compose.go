package docker

import (
	"os/exec"
)

type Compose struct {
	ComposeFile string
	runner
}

func (c Compose) run(args ...string) error {
	a := []string{}
	if v := c.ComposeFile; v != "" {
		a = append(a, "--file", v)
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
