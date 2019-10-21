package docker

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/duckbrain/beluga/internal/portainer"
)

// TODO: Portainer has an API endpoint to deploy a stack without relaying the docker proxy.
// https://app.swaggerhub.com/apis-docs/deviantony/Portainer/1.22.0#/stacks/StackCreate
// For this to be used, portainer.Client should implement the Deployer interface.
type portainerRun struct {
	port   int
	client *portainer.Client
}

func newPortainer(u *url.URL) *portainerRun {
	client := &portainer.Client{
		Client: http.Client{
			Timeout: time.Minute,
		},
		DSN: u,
	}
	return &portainerRun{client: client}
}

func (cmd *portainerRun) listen() error {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}
	cmd.port = ln.Addr().(*net.TCPAddr).Port
	go http.Serve(ln, cmd.client)
	return nil
}

func (cmd *portainerRun) run(c *exec.Cmd) error {
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if cmd.port == 0 {
		err := cmd.client.Authenticate()
		if err != nil {
			return err
		}
		err = cmd.listen()
		if err != nil {
			return err
		}
	}
	c.Env = []string{fmt.Sprintf("DOCKER_HOST=http://localhost:%v", cmd.port)}
	return c.Run()
}
