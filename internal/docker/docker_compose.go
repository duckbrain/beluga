package docker

import "os/exec"

type Deployer interface {
	Deploy(composeFile string) error
	Teardown(composeFile string) error
}
type compose struct{}

var Compose Deployer = compose{}

func (compose) Deploy(composeFile string) error {
	return exec.Command(
		"docker-compose", "up", "-d",
		"--no-build",
		"-f", composeFile,
	).Run()
}

func (compose) Teardown(composeFile string) error {
	return exec.Command(
		"docker-compose", "down",
		"-f", composeFile,
	).Run()
}
