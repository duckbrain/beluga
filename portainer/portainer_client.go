package portainer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var schemeMap = map[string]string{
	"http":               "http",
	"portainer-insecure": "http",
	"https":              "https",
	"portainer":          "https",
}

func New(dsn string, opts interface{}) (*Client, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}
	client := &Client{
		Credentials: u.User,
		Client:      http.DefaultClient,
		Logger:      &logrus.Logger{},
	}
	if scheme, ok := schemeMap[u.Scheme]; ok {
		client.BaseURL = fmt.Sprintf("%v://%v", scheme, u.Host)
	} else {
		return nil, errors.Errorf("unknown scheme %v", u.Scheme)
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	if opts != nil {
		err = decoder.Decode(opts, map[string][]string(u.Query()))
		if err != nil {
			return nil, errors.Wrap(err, "decode filters")
		}
	}

	return client, nil
}

type Client struct {
	BaseURL     string
	Client      *http.Client
	Credentials *url.Userinfo
	JWT         string
	Logger      logrus.StdLogger
}

type Endpoint struct {
	ID   int64  `json:"Id"`
	Name string `json:"Name"`
}

type Endpoints []Endpoint

type StackType int64

const (
	Swarm   StackType = 1
	Compose StackType = 2
)

type Env map[string]string

func (env Env) MarshalJSON() ([]byte, error) {
	type entry struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	var values []entry
	for name, value := range env {
		values = append(values, entry{name, value})
	}
	return json.Marshal(values)
}

func (env Env) UnmarshalJSON(data []byte) error {
	type entry struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	var values []entry
	err := json.Unmarshal(data, &values)
	if err != nil {
		return err
	}
	for _, entry := range values {
		env[entry.Name] = entry.Value
	}
	return nil
}

type Stack struct {
	ID         int64     `json:"Id"`
	Name       string    `json:"Name"`
	Type       StackType `json:"Type"`
	EndpointID int64     `json:"EndpointId"`
	Env        Env       `json:"Env"`
}

type Stacks []Stack

type Error struct {
	Message string `json:"err"`
}

func (e Error) Error() string {
	return e.Message
}

func (c *Client) do(method, path string, params, body, response interface{}) error {
	if params != nil {
		urlParams := map[string][]string{}
		err := schema.NewEncoder().Encode(params, urlParams)
		if err != nil {
			return errors.Wrap(err, "encode URL parameters")
		}
		path += "?" + url.Values(urlParams).Encode()
	}
	reqJSON, err := json.Marshal(&body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, c.BaseURL+path, bytes.NewReader(reqJSON))
	if err != nil {
		return err
	}
	c.Logger.Printf("%v %v", req.Method, req.URL.String())
	if c.JWT != "" {
		req.Header.Set("Authorization", "Bearer "+c.JWT)
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	buffer := &bytes.Buffer{}
	bodyDecoder := json.NewDecoder(io.TeeReader(resp.Body, buffer))
	if resp.StatusCode == 200 {
		err = errors.Wrap(bodyDecoder.Decode(response), "parsing response")
		c.Logger.Printf("response: %v", buffer.String())
		return err
	}
	errorMessage := Error{}
	err = bodyDecoder.Decode(&errorMessage)
	if err != nil {
		err = errors.Wrap(err, "parsing error message")
		c.Logger.Printf("response %v: %v", resp.Status, buffer.String())
		return err
	}
	return errorMessage
}

// Authenticate logs in with the user credentials, and obtains a JWT to use for
// other calls. This must be called before any other calls.
func (c *Client) Authenticate(u *url.Userinfo) (string, error) {
	if u == nil {
		u = c.Credentials
	}
	if u == nil {
		return "", errors.New("no credentials provided")
	}
	var reqBody struct {
		Username string
		Password string
	}
	var resBody struct{ JWT string }
	reqBody.Username = u.Username()
	if pasword, ok := u.Password(); ok {
		reqBody.Password = pasword
	} else {
		return "", errors.New("no password provided")
	}
	err := c.do("POST", "/api/auth", nil, reqBody, &resBody)
	return resBody.JWT, err
}

func (c *Client) Endpoint(id int64) (e Endpoint, err error) {
	err = c.do("GET", fmt.Sprintf("/api/endpoints/%v", id), nil, nil, &e)
	return
}

type EndpointsFilter struct {
	GroupID int64 `schema:"groupId"`
}

func (c *Client) Endpoints(filter EndpointsFilter) (e Endpoints, err error) {
	err = c.do("GET", "/api/endpoints", filter, nil, &e)
	return
}

type StacksFilter struct {
	EndpointID int64  `json:"EndpointId,omitempty"`
	SwarmID    string `json:"SwarmId,omitempty"`
}

func (c *Client) Stacks(filters StacksFilter) (stacks Stacks, err error) {
	var params interface{}
	data, err := json.Marshal(filters)
	if err != nil {
		return nil, errors.Wrap(err, "encode filters")
	}
	if string(data) != "{}" {
		params = struct{ Filters string }{Filters: string(data)}
	}
	err = c.do("GET", "/api/stacks", params, nil, &stacks)
	return
}

// https://app.swaggerhub.com/apis-docs/deviantony/Portainer/1.22.0#/stacks/StackCreate
func (c *Client) NewStack(s Stack, composeFileContents string) (Stack, error) {
	params := struct {
		Type       int64  `schema:"type"`
		Method     string `schema:"method"`
		EndpointID int64  `schema:"endpointId"`
	}{
		Type:       int64(s.Type),
		Method:     "string",
		EndpointID: s.EndpointID,
	}

	body := struct {
		Name             string
		SwarmID          string
		StackFileContent string
		Env              Env
	}{
		Name:             s.Name,
		StackFileContent: composeFileContents,
		Env:              s.Env,
	}

	newStack := Stack{}

	err := c.do("POST", "/api/stacks", params, body, &newStack)

	return newStack, err
}
