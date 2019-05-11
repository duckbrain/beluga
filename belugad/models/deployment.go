package models

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/duckbrain/beluga/beluga"
)

type Deployment = beluga.Deployment

type Deployer interface {
	Deploy(stackName string, d Deployment) error
	Teardown(stackName string, d Deployment) error
}

type Compose struct {
	Cache    string
	Executor string
	Command  string
}

func (c Compose) write(stackName, fileName, contents string) error {
	return ioutil.WriteFile(filepath.Join(c.Cache, stackName, fileName), []byte(contents), os.ModePerm)
}

func (c Compose) run(stackName string, d Deployment, args ...string) error {
	if err := os.MkdirAll(filepath.Join(c.Cache, stackName), os.ModePerm); err != nil {
		return err
	}
	if err := c.write(stackName, "docker-compose.yaml", d.ComposeFile); err != nil {
		return err
	}
	if err := c.write(stackName, ".env", d.Envfile()); err != nil {
		return err
	}
	cmd := exec.Command(c.Executor, args...)
	cmd.Dir = filepath.Join(c.Cache, stackName)
	return cmd.Run()
}

func (c Compose) Deploy(stackName string, d Deployment) error {
	return c.run(stackName, d, c.Command, "up", "-d")
}

func (c Compose) Teardown(stackName string, d Deployment) error {
	return c.run(stackName, d, c.Command, "down")
}
