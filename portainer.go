package beluga

import (
	"github.com/duckbrain/beluga/portainer"
)

type PortainerDeploy struct {
	Client     *portainer.Client
	StackType  string
	EndpointID int64
	GroupID    int64
}

func (c *PortainerDeploy) Deploy(opts DeployOpts) error {
	panic("not implemented")
}

func (c *PortainerDeploy) Teardown(opts DeployOpts) error {
	panic("not implemented")
}
