package beluga

import (
	"net/url"
)

type DeployMode int

const (
	ComposeMode DeployMode = 2
	SwarmMode   DeployMode = 1
)

type Deployer interface {
	Deploy(opts DeployOpts) error
	Teardown(opts DeployOpts) error
}

func (env Environment) Deployer() Deployer {
	host := env[varDeployDockerHost]
	u, _ := url.Parse(host)
	switch u.Scheme {
	case "portainer":
		return &Portainer{DSN: u}
	case "ssh":
		panic("SSH not implemented")
	default:
		return Docker(host)
	}
}
