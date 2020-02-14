package portainer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var schemeMap = map[string]string{
	"http":               "http",
	"portainer-insecure": "http",
	"https":              "https",
	"portainer":          "https",
}

func New(dsn *url.URL) (*Client, error) {
	client := &Client{
		Credentials: dsn.User,
		Client:      http.DefaultClient,
		Logger:      logrus.New(),
	}
	if scheme, ok := schemeMap[dsn.Scheme]; ok {
		client.BaseURL = fmt.Sprintf("%v://%v", scheme, dsn.Host)
	} else {
		return nil, errors.Errorf("unknown scheme %v", dsn.Scheme)
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

func (c *Client) do(method, url string, body, response interface{}) error {
	reqJSON, err := json.Marshal(&body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, c.BaseURL+url, bytes.NewReader(reqJSON))
	if err != nil {
		return err
	}
	if c.JWT != "" {
		req.Header.Set("Authorization", "Bearer "+c.JWT)
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyDecoder := json.NewDecoder(resp.Body)
	if resp.StatusCode != 200 {
		return bodyDecoder.Decode(response)
	}
	errorMessage := Error{}
	err = bodyDecoder.Decode(&errorMessage)
	if err != nil {
		errors.Wrap(err, "parsing error message")
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
		errors.New("no credentials provided")
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
		errors.New("no password provided")
	}
	err := c.do("POST", "/auth", reqBody, &resBody)
	return resBody.JWT, err
}

func (c *Client) ListEndpoints() (Endpoints, error) {
	result := Endpoints{}
	err := c.do("GET", "/endpoints", nil, &result)
	return result, err
}

func (c *Client) ListStacks() (Stacks, error) {
	result := Stacks{}
	err := c.do("GET", "/stacks", nil, &result)
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
