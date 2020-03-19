package beluga

import (
	"context"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Runner struct {
	Env    Environment
	Logger logrus.StdLogger
}

func New() Runner {
	return Runner{
		Env:    Env(),
		Logger: logrus.New(),
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

	return composeTemplate(string(out), r.Env.ComposeTemplate(), r.Env)
}

func (r Runner) StackName() string {
	return r.Env[varStackName]
}

func (r Runner) Deploy(ctx context.Context) error {
	return r.deployer().Deploy(ctx, r)
}

func (r Runner) Teardown(ctx context.Context) error {
	return r.deployer().Deploy(ctx, r)
}
