package docker

import (
	"net/url"
	"os"
	"os/exec"

	"github.com/duckbrain/beluga/internal/lib"
)

type Docker struct {
	runner
}

type runner interface {
	run(*exec.Cmd) error
}

var Local = Docker{runner: localRun{}}

func New(host string) Docker {
	if host == "" {
		return Local
	}

	u, err := url.Parse(host)
	if err != nil {
		return Local
	}

	switch u.Scheme {
	case "portainer":
		return Docker{runner: newPortainer(u)}
	default:
		return Local
	}
}

type localRun struct{}

func (localRun) run(c *exec.Cmd) error {
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func (d Docker) Build(context, dockerfile, tag string) error {
	return d.run(exec.Command(
		"docker", "build",
		context,
		"-f", dockerfile,
		"-t", tag,
	))
}

func (d Docker) Tag(src, dst string) error {
	return d.run(exec.Command("docker", "tag", src, dst))
}

func (d Docker) Push(tag string) error {
	return d.run(exec.Command("docker", "push", tag))
}

func (d Docker) Login(hostname, username, password string) error {
	return d.run(exec.Command(
		"docker", "login",
		hostname,
		"-u", username,
		"-p", password,
	))
}

func (d Docker) Compose(env lib.Environment) Compose {
	return Compose{env: env, runner: d.runner}
}
