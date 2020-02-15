package portainer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gobwas/glob"
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

func New(dsn string) (*Client, error) {
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

	err = decoder.Decode(&client.Filters, map[string][]string(u.Query()))
	if err != nil {
		return nil, errors.Wrap(err, "decode filters")
	}

	return client, nil
}

type Client struct {
	BaseURL     string
	Client      *http.Client
	Credentials *url.Userinfo
	JWT         string
	Logger      logrus.StdLogger
	Filters     struct {
		ID      *int64 `schema:"id"`
		GroupID *int64 `schema:"groupId"`
	}
}

type Endpoint struct {
	ID   int64  `json:"Id"`
	Name string `json:"Name"`
}

type Endpoints []Endpoint

func (list Endpoints) Filter(s string) Endpoints {
	newList := Endpoints{}

	if id, err := strconv.ParseInt(s, 10, 64); err == nil {
		for _, endpoint := range list {
			if endpoint.ID == int64(id) {
				newList = append(newList, endpoint)
			}
		}

	} else if g, err := glob.Compile(s); err == nil {
		for _, endpoint := range list {
			if g.Match(endpoint.Name) {
				newList = append(newList, endpoint)
			}
		}
	} else {
		for _, endpoint := range list {
			if endpoint.Name == s {
				newList = append(newList, endpoint)
				break
			}
		}
	}

	return newList
}

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

type Stack struct {
	ID         string `json:"Id"`
	Name       string `json:"Name"`
	Type       StackType
	EndpointID int64
	Env        Env
}

type Stacks []Stacks

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

func (c *Client) ListEndpoints() (Endpoints, error) {
	result := Endpoints{}
	err := c.do("GET", "/api/endpoints", c.Filters, nil, &result)
	return result, err
}

func (c *Client) ListStacks() (Stacks, error) {
	result := Stacks{}
	err := c.do("GET", "/api/stacks", nil, nil, &result)
	return result, err
}

// TODO https://app.swaggerhub.com/apis-docs/deviantony/Portainer/1.22.0#/stacks/StackCreate
func (c *Client) CreateStack(s Stack) error {
	panic("TODO")
	// 	queryParameters := url.Values{
	// 		"type":       {fmt.Sprint(opts.DeployMode())},
	// 		"method":     {"string"},
	// 		"endpointID": {strings.TrimPrefix(c.DSN.Path, "/")},
	// 	}
	// 	composeFileContents, err := opts.ComposeFileContents()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	var requestBody struct {
	// 		Name             string
	// 		StackFileContent string
	// 	}
	// 	requestBody.Name = opts.StackName()
	// 	requestBody.StackFileContent = composeFileContents
	// 	reqJSON, err := json.Marshal(&requestBody)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	u := c.path("/stacks") + "?" + queryParameters.Encode()

	// 	c.Logger.Printf("POST %v %v", u, string(reqJSON))
	// 	resp, err := c.Client.Post(u, "application/json", bytes.NewReader(reqJSON))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	defer resp.Body.Close()
	// 	if resp.StatusCode == 200 {
	// 		return nil
	// 	}
	// 	return nil
}
