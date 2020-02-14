package beluga

import (
	"github.com/duckbrain/beluga/portainer"
)

type PortainerDeploy struct {
	Client *portainer.Client
}

func (c *PortainerDeploy) Deploy(opts DeployOpts) error {
	panic("not implemented")
}

func (c *PortainerDeploy) Teardown(opts DeployOpts) error {
	panic("not implemented")
}
