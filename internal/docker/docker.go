package docker

import (
	"net/url"
	"os"
	"os/exec"

	"github.com/duckbrain/beluga/internal/lib"
)

type Docker struct {
	commander
}

type commander interface {
	cmd(s string, args ...string) *exec.Cmd
}

var Local = Docker{commander: localCmd{}}

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
		return Docker{commander: newPortainer(u)}
	default:
		return Local
	}
}

type localCmd struct{}

func (localCmd) cmd(s string, args ...string) *exec.Cmd {
	c := exec.Command(s, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c
}

func (d Docker) Build(context, dockerfile, tag string) error {
	c := d.cmd(
		"docker", "build",
		context,
		"-f", dockerfile,
		"-t", tag,
	)
	return c.Run()
}

func (d Docker) Tag(src, dst string) error {
	return d.cmd("docker", "tag", src, dst).Run()
}

func (d Docker) Push(tag ...string) error {
	args := []string{"push"}
	args = append(args, tag...)
	return d.cmd("docker", args...).Run()
}

func (d Docker) Login(hostname, username, password string) error {
	return d.cmd(
		"docker", "login",
		hostname,
		"-u", username,
		"-p", password,
	).Run()
}

func (d Docker) Compose(env lib.Environment) Compose {
	return Compose{env: env, commander: d.commander}
}
