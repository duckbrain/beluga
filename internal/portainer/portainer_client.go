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

var statusMap = map[StatusError]string{
	403: "Unauthorized",
	404: "Not found",
}

type StatusError int64

func (e StatusError) Error() string {
	msg, ok := statusMap[e]
	if ok {
		return msg
	}
	return fmt.Sprintf("status code %v", int64(e))
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

func (e Endpoints) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

type StackType int64

const (
	Swarm   StackType = 1
	Compose StackType = 2
)

var stackTypeMap = map[string]StackType{
	"swarm":   Swarm,
	"compose": Compose,
}

func ParseStackType(s string) (StackType, error) {
	t, ok := stackTypeMap[s]
	if !ok {
		return 0, errors.New("unknown stack type")
	}
	return t, nil
}

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

func (env *Env) UnmarshalJSON(data []byte) error {
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
		if *env == nil {
			*env = make(Env)
		}
		(*env)[entry.Name] = entry.Value
	}
	return nil
}

type Stack struct {
	ID         int64  `json:"Id"`
	Name       string `json:"Name"`
	EndpointID int64  `json:"EndpointId"`
	Env        Env    `json:"Env"`
}

type Stacks []Stack

type Error struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

func (e Error) Error() string {
	return fmt.Sprintf("server error: %v; %v", e.Message, e.Details)
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
	if resp.StatusCode >= 200 && resp.StatusCode <= 400 {
		if response == nil {
			return nil
		}
		err = errors.Wrap(bodyDecoder.Decode(response), "parsing response")
		c.Logger.Printf("response: %v", buffer.String())
		return err
	}
	if _, ok := statusMap[StatusError(resp.StatusCode)]; ok {
		return StatusError(resp.StatusCode)
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

func (c *Client) SwarmID(endpointID int64) (string, error) {
	var response struct {
		ID string
	}
	err := c.do("GET", fmt.Sprintf("/api/endpoints/%v/docker/swarm", endpointID), nil, nil, &response)
	return response.ID, err
}

// https://app.swaggerhub.com/apis-docs/deviantony/Portainer/1.22.0#/stacks/StackCreate
func (c *Client) NewStack(s Stack, composeFileContents string) (newStack Stack, err error) {
	params := struct {
		Type       StackType `schema:"type"`
		Method     string    `schema:"method"`
		EndpointID int64     `schema:"endpointId"`
	}{
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

	// When deploying in swarm mode, we must include the swarm ID of the
	// entrypoint. We can get this by querying the docker socket.
	// See: https://gist.github.com/deviantony/77026d402366b4b43fa5918d41bc42f8#manage-docker-stacks-in-a-swarm-environment
	body.SwarmID, err = c.SwarmID(s.EndpointID)
	if err != nil {
		err = errors.Wrap(err, "fetch swarm ID")
		return
	}
	if body.SwarmID == "" {
		params.Type = Compose
	} else {
		params.Type = Swarm
	}

	err = c.do("POST", "/api/stacks", params, body, &newStack)
	return
}

func (c *Client) UpdateStack(s Stack, composeFileContents string, prune bool) (Stack, error) {
	var params interface{}
	if s.EndpointID > 0 {
		params = struct {
			EndpointID int64 `schema:"endpointId"`
		}{EndpointID: s.EndpointID}
	}

	body := struct {
		StackFileContent string
		Env              Env
		Prune            bool
	}{
		StackFileContent: composeFileContents,
		Env:              s.Env,
		Prune:            prune,
	}

	updatedStack := Stack{}

	err := c.do("PUT", fmt.Sprintf("/api/stacks/%v", s.ID), params, body, &updatedStack)

	return updatedStack, err
}

func (c *Client) RemoveStack(id int64) error {
	return c.do("DELETE", fmt.Sprintf("/api/stacks/%v", id), nil, nil, nil)
}
