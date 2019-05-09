package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type Deployment struct {
	Env         map[string]string `json:"env"`
	ComposeFile string            `json:"compose_file"`
}

type Deployer interface {
	Deploy(stackName string, d Deployment) error
	Teardown(stackName string, d Deployment) error
}

type Compose struct {
	Cache    string
	Executor string
	Command  string
}

func (d Deployment) Envfile() string {
	lines := make([]string, 0, len(d.Env))
	for k, v := range d.Env {
		lines = append(lines, fmt.Sprintf(`%s=%s`, k, v))
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n")

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
