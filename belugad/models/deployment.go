package models

type Deployment struct {
	Env        map[string]string `json:"env"`
	Dockerfile string            `json:"dockerfile"`
}

type Deployer interface {
	Deploy(stackName string, d Deployment) error
	Teardown(stackName string) error
}
