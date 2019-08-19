package docker

type Deployer interface {
	Deploy(composeFile string) error
	Teardown(composeFile string) error
}
type compose struct{}

var Compose Deployer = compose{}

func (compose) Deploy(composeFile string) error {
	return run(
		"docker-compose", "up", "-d",
		"--no-build",
		"-f", composeFile,
	)
}

func (compose) Teardown(composeFile string) error {
	return run(
		"docker-compose", "down",
		"-f", composeFile,
	)
}
