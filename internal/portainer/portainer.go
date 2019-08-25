package portainer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

type Client struct {
	DSN    url.URL
	Client http.Client

	jwt string
}

func (c Client) path(s string) string {
	u := c.DSN
	u.Path = s
	u.User = nil
	u.Fragment = ""
	u.RawPath = s
	u.RawQuery = ""
	return u.String()
}

func (c *Client) Authenticate() error {
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

func (c *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := c.path(path.Join("/endpoints", c.DSN.Path, "docker"), r.URL.EscapedPath()))
	r.URL url.Parse(p)
}
