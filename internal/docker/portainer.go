package docker

import (
	"net/http"
	"net/url"
	"os/exec"
	"time"

	"github.com/duckbrain/beluga/internal/portainer"
)

type portainerCmd struct{ client portainer.Client }

func newPortainer(u *url.URL) portainerCmd {
	client := portainer.Client{
		Client: http.Client{
			Timeout: time.Minute,
		},
		DSN: u,
	}
	return portainerCmd{client: client}
}

func (portainerCmd) cmd(s string, args ...string) *exec.Cmd {
	c := localCmd{}.cmd(s, args...)
	return c
}
