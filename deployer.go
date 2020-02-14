package beluga

import (
	"net/url"

	"github.com/duckbrain/beluga/portainer"
	"github.com/sirupsen/logrus"
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

func (env Environment) Logger() logrus.StdLogger {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	return logger
}

func (env Environment) Deployer() Deployer {
	logger := env.Logger()
	host := env[varDeployDockerHost]
	u, _ := url.Parse(host)
	switch u.Scheme {
	case "portainer", "portainer-insecure":
		logger.Printf("Portainer url: %v", u)
		client, err := portainer.New(u)
		if err != nil {
			panic(err)
		}
		client.Logger = logger
		return &PortainerDeploy{Client: client}
	case "ssh":
		panic("SSH not implemented")
	default:
		return Docker(host)
	}
}
