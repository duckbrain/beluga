package beluga

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type Portainer struct {
	// DSN is a parsed URL that represents the domain specific language for the
	// connection. The Scheme must be "portainer", (this is changed to HTTPs)
	// for requests. The host and optional port are used literally, the username
	// and password are used to authenticate, and the path is the endpoint to use.
	DSN *url.URL

	// Client is required to be set. It must be set to the HTTP client to use to
	// make requests.
	Client http.Client

	jwt string
}

func (c Portainer) path(s ...string) string {
	u := *c.DSN
	u.Path = path.Join(s...)
	u.User = nil
	u.Fragment = ""
	u.RawQuery = ""
	if c.flag("http") {
		u.Scheme = "http"
	} else {
		u.Scheme = "https"
	}
	return u.String()
}

// flag returns true if flag s is present on the DSN
func (c Portainer) flag(s string) bool {
	_, ok := c.DSN.Query()[s]
	return ok
}

// Authenticate logs in with the user credentials, and obtains a JWT to use for
// other calls. This must be called before any other calls.
func (c *Portainer) Authenticate() error {
	if c.DSN.User == nil {
		return errors.New("no user info in DNS")
	}
	reqBody := struct {
		Username string
		Password string
	}{}
	reqBody.Username = c.DSN.User.Username()
	reqBody.Password, _ = c.DSN.User.Password()
	reqJSON, err := json.Marshal(&reqBody)
	if err != nil {
		return err
	}
	res, err := c.Client.Post(c.path("/auth"), "application/json", bytes.NewReader(reqJSON))
	if err != nil {
		return err
	}
	var resBody struct{ JWT string }
	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return err
	}
	c.jwt = resBody.JWT
	return nil
}

func (c *Portainer) serveHTTPError(err error, w http.ResponseWriter, r *http.Request) {
	r.Body.Close()
	w.WriteHeader(503)
	_, _ = w.Write([]byte(fmt.Sprintf("HTTP proxy error: %v", err.Error())))
}

// ServeHTTP implements the http.Handler interface to allow listening on and
// proxying requests from a docker client to the deamon, authenticated through
// the Portainer instance
func (c *Portainer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := c.path("/endpoints", c.DSN.Path, "docker", r.URL.Path)
	req, err := http.NewRequest(r.Method, p, r.Body)
	if err != nil {
		c.serveHTTPError(err, w, r)
		return
	}
	r.Body.Close()
	res, err := c.Client.Do(req)
	if err != nil {
		c.serveHTTPError(err, w, r)
		return
	}
	for k, vs := range res.Header {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(res.StatusCode)
	_, err = io.Copy(w, res.Body)
	if err != nil {
		c.serveHTTPError(err, w, r)
		return
	}
	res.Body.Close()
}

// TODO https://app.swaggerhub.com/apis-docs/deviantony/Portainer/1.22.0#/stacks/StackCreate
func (c *Portainer) Deploy(opts DeployOpts) error {
	queryParameters := url.Values{
		"type":       {fmt.Sprint(opts.DeployMode())},
		"method":     {"string"},
		"endpointID": {c.DSN.Path},
	}
	composeFileContents, err := opts.ComposeFileContents()
	if err != nil {
		return err
	}
	var requestBody struct {
		Name             string
		StackFileContent string
	}
	requestBody.Name = opts.StackName()
	requestBody.StackFileContent = composeFileContents
	reqJSON, err := json.Marshal(&requestBody)
	if err != nil {
		return err
	}
	_, err = c.Client.Post(
		c.path("/stacks")+"?"+queryParameters.Encode(),
		"application/json",
		bytes.NewReader(reqJSON))
	return err
}

// TODO https://app.swaggerhub.com/apis-docs/deviantony/Portainer/1.22.0#/stacks/StackDelete
func (c *Portainer) Teardown(opts DeployOpts) error {
	panic("not implemented")
}
