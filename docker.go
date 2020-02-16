package beluga

import (
	"io/ioutil"
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

type BuildInfo interface {
	BuildContext() string
	Dockerfile() string
	DockerImage() string
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

type DeployOpts interface {
	ComposeFileContents() (string, error)
	StackName() string
	DeployMode() DeployMode
}

func (env Environment) ComposeFileContents() (string, error) {
	cmd := exec.Command("docker-compose", "config")
	env["COMPOSE_FILE"] = env[varDockerComposeFile]
	cmd.Env, _ = env.Format(GoEnvFormat, true)
	out, err := cmd.Output()
	return string(out), err
}

func writeComposeFile(opts DeployOpts) (filename string, err error) {
	contents, err := opts.ComposeFileContents()
	file, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	_, err = file.WriteString(contents)
	if err != nil {
		return "", err
	}
	return file.Name(), file.Close()
}

func (d Docker) Deploy(opts DeployOpts) error {
	stackName := opts.StackName()
	composeFile, err := writeComposeFile(opts)
	if err != nil {
		return err
	}
	defer os.Remove(composeFile)

	var cmd *exec.Cmd
	switch opts.DeployMode() {
	case ComposeMode:
		cmd = exec.Command("docker-compose",
			"--file", composeFile,
			"--project-name", stackName,
			"up",
			"--detach", "--no-build")
	case SwarmMode:
		cmd = exec.Command("docker", "stack", "deploy",
			"--compose-file", composeFile,
			"--prune",
			"--with-registry-auth",
			stackName)
	}
	return d.run(cmd)
}

func (d Docker) Teardown(opts DeployOpts) error {
	stackName := opts.StackName()
	var cmd *exec.Cmd
	switch opts.DeployMode() {
	case ComposeMode:
		composeFile, err := writeComposeFile(opts)
		if err != nil {
			return err
		}
		defer os.Remove(composeFile)
		cmd = exec.Command("docker-compose",
			"--file", composeFile,
			"--project-name", stackName,
			"down",
			"--volumes", "--remove-orphans")
	case SwarmMode:
		cmd = exec.Command("docker", "stack", "rm", stackName)
	}
	return d.run(cmd)
}
