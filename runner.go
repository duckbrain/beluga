package beluga

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Runner struct {
	Env    Environment
	Logger logrus.StdLogger
	DryRun bool
}

func New() *Runner {
	logger := logrus.New()

	return &Runner{
		Env:    Env(logger),
		Logger: logger,
	}
}

// Exec runs a command in the Beluga context with stderr and stdout bound to the parent process
func (r *Runner) Exec(c *exec.Cmd) error {
	r.Logger.Println(c.String())
	if r.DryRun {
		return nil
	}
	var err error
	if c.Stdin == nil {
		c.Stdin = os.Stdin
	}
	if c.Stdout == nil {
		c.Stdout = os.Stdout
	}
	if c.Stderr == nil {
		c.Stderr = os.Stderr
	}
	c.Env, err = r.Env.Format(GoEnvFormat, r.Env.Keys())
	if err != nil {
		return errors.Wrap(err, "generate environment list")
	}
	return c.Run()
}

func (r *Runner) ComposeFile(ctx context.Context) (string, error) {
	d := r.docker()
	out, err := d.ComposeConfig(ctx)
	if err != nil {
		return "", err
	}

	return composeTemplate(string(out), r.Env.ComposeTemplate(), r.Env)
}

func (r *Runner) StackName() string {
	return r.Env[varStackName]
}

func (r *Runner) Deploy(ctx context.Context) error {
	return r.deployer().Deploy(ctx, r)
}

func (r *Runner) Teardown(ctx context.Context) error {
	return r.deployer().Deploy(ctx, r)
}

type BuildOpts struct {
	Push bool
}

func (r *Runner) Build(ctx context.Context, opts BuildOpts) error {
	e := r.Env
	d := r.docker()

	if e.RegistryUsername() != "" {
		err := d.Login(ctx, e.Registry(), e.RegistryUsername(), e.RegistryPassword())
		if err != nil {
			return err
		}
	}

	builtImage := ""
	images := strings.Fields(e.Images())

	for _, image := range images {
		if builtImage == "" {
			err := d.Build(ctx, e.Context(), e.Dockerfile(), image)
			if err != nil {
				return err
			}
			builtImage = image
		} else {
			err := d.Tag(ctx, builtImage, image)
			if err != nil {
				return err
			}
		}
	}

	if opts.Push {
		for _, image := range images {
			err := d.Push(ctx, image)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
