package lib

import (
	"net/url"

	"github.com/duckbrain/beluga/internal/portainer"

	"github.com/duckbrain/beluga/internal/docker"
)

type Deployer interface {
	Deploy(composeFile string, env map[string]string) error
	Teardown(composeFile string) error
}

func (env Environment) Deployer() Deployer {
	host := env.DeployDockerHost()
	u, _ := url.Parse(host)
	switch u.Scheme {
	case "portainer":
		return &portainer.Client{DSN: u}
	default:
		return docker.New(host).Compose()
	}
}
