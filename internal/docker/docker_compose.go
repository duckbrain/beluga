package docker

import (
	"os/exec"
)

func cliEnv(in map[string]string) []string {
	out := make([]string, 0, len(in))
	for key, value := range in {
		out = append(out, key+"="+value)
	}
	return out
}

type Compose struct {
	ComposeFile string
	runner
}

func (c Compose) run(composeFile string, env map[string]string, args ...string) error {
	a := []string{}
	if v := c.ComposeFile; v != "" {
		a = append(a, "--file", v)
	}
	a = append(a, args...)
	cmd := exec.Command("docker-compose", a...)
	cmd.Env = cliEnv(env)
	return c.runner.run(cmd)
}

func (c Compose) Deploy(composeFile string, env map[string]string) error {
	return c.run(composeFile, env, "up", "--detach", "--no-build")
}

func (c Compose) Teardown(composeFile string) error {
	return c.run(composeFile, nil, "down")
}
