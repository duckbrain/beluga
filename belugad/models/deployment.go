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
	Deploy(domain string, d Deployment) error
	Teardown(domain string, d Deployment) error
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

func (c Compose) write(domain, fileName, contents string) error {
	return ioutil.WriteFile(filepath.Join(c.Cache, domain, fileName), []byte(contents), os.ModePerm)
}

func (c Compose) run(domain string, d Deployment, args ...string) error {
	if err := os.MkdirAll(filepath.Join(c.Cache, domain), os.ModePerm); err != nil {
		return err
	}
	if err := c.write(domain, "docker-compose.yaml", d.ComposeFile); err != nil {
		return err
	}
	if err := c.write(domain, ".env", d.Envfile()); err != nil {
		return err
	}
	cmd := exec.Command(c.Executor, args...)
	cmd.Dir = filepath.Join(c.Cache, domain)
	return cmd.Run()
}

func (c Compose) Deploy(domain string, d Deployment) error {
	return c.run(domain, d, c.Command, "up", "-d")
}

func (c Compose) Teardown(domain string, d Deployment) error {
	return c.run(domain, d, c.Command, "down")
}
