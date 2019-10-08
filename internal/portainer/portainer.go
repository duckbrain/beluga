package portainer

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

type Client struct {
	DSN    *url.URL
	Client http.Client

	jwt string
}

func (c Client) path(s ...string) string {
	u := *c.DSN
	u.Path = path.Join(s...)
	u.User = nil
	u.Fragment = ""
	u.RawQuery = ""
	u.Scheme = "https"
	return u.String()
}

func (c *Client) Authenticate() error {
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

func (c *Client) serveHTTPError(err error, w http.ResponseWriter, r *http.Request) {
	r.Body.Close()
	w.WriteHeader(503)
	w.Write([]byte(fmt.Sprintf("HTTP proxy error: %v", err.Error())))
}

func (c *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
