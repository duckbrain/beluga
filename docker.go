package beluga

import (
	"os"
	"os/exec"
)

type Docker string

func (d Docker) run(c *exec.Cmd) error {
	if c.Stdin == nil {
		c.Stdin = os.Stdin
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env = []string{"DOCKER_HOST=" + string(d)}
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

func (d Docker) Deploy(opts DeployOpts) error {
	composeFile, err := opts.writeComposeFile()
	if err != nil {
		return err
	}
	defer os.Remove(composeFile)

	var cmd *exec.Cmd
	switch opts.Mode {
	case ComposeMode:
		cmd = exec.Command("docker-compose",
			"--file", composeFile,
			"--project-name", opts.StackName,
			"up",
			"--detach", "--no-build")
	case SwarmMode:
		cmd = exec.Command("docker", "stack", "deploy",
			"--compose-file", composeFile,
			"--prune",
			"--with-registry-auth",
			opts.StackName)
	}
	return d.run(cmd)
}

func (d Docker) Teardown(opts DeployOpts) error {
	var cmd *exec.Cmd
	switch opts.Mode {
	case ComposeMode:
		composeFile, err := opts.writeComposeFile()
		if err != nil {
			return err
		}
		defer os.Remove(composeFile)
		cmd = exec.Command("docker-compose",
			"--file", composeFile,
			"--project-name", opts.StackName,
			"down",
			"--volumes", "--remove-orphans")
	case SwarmMode:
		cmd = exec.Command("docker", "stack", "rm", opts.StackName)
	}
	return d.run(cmd)
}
