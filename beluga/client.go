// Package beluga provides shared functionality for the server and CLI.
package beluga

import (
	"errors"
	"net/http"
)

type Client struct {
	DSN    string
	Client *http.Client
}

type Deployment struct {
	Env         map[string]string `json:"env"`
	ComposeFile string            `json:"compose_file"`
}

func (c Client) http(method, stackName string, d *Deployment) error {
	return errors.New("Not implemented")
}

func (c Client) Deploy(stackName string, d Deployment) error {
	return c.http("PUT", stackName, &d)
}
func (c Client) Teardown(stackName string) error {
	return c.http("DELETE", stackName, nil)
}
