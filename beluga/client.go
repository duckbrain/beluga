// Package beluga provides go Client library for the beluga server.
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

func (c Client) http(method, domain string, d *Deployment) error {
	return errors.New("Not implemented")
}

func (c Client) Deploy(domain string, d Deployment) error {
	return c.http("PUT", domain, &d)
}
func (c Client) Teardown(domain string) error {
	return c.http("DELETE", domain, nil)
}
