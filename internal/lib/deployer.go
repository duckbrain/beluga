package lib

import (
	"github.com/duckbrain/beluga/internal/docker"
)

type Deployer interface {
	Deploy() error
	Teardown() error
}


func (env Environment) Deployer() Deployer {
	return docker.New(env.DeployDockerHost()).Compose(env.DockerComposeFile())
}