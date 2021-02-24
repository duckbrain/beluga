package beluga

import (
	"bytes"
	"context"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/duckbrain/beluga/internal/portainer"
	"github.com/pkg/errors"
)

type DeployMode int

const (
	ComposeMode DeployMode = 2
	SwarmMode   DeployMode = 1
)

type dummyTransport struct{}

func (d dummyTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, nil
}

func (r Runner) deployer() Deployer {
	host := r.Env.DeployDSN()

	if i := strings.Index(host, ":"); i > 0 {
		switch host[:i] {
		case "portainer", "portainer-insecure":
			deployer := &portainerDeploy{}
			client, err := portainer.New(host, deployer)
			if err != nil {
				panic(err)
			}
			if r.DryRun {
				client.Client.Transport = dummyTransport{}
			}
			client.Logger = r.Logger
			deployer.Client = client
			return deployer
		case "ssh":
			panic("SSH not implemented")
		}
	}
	return r.docker()
}

func (r Runner) docker() dockerRunner {
	return dockerRunner{Exec: r.Exec}
}

type Deployer interface {
	Deploy(ctx context.Context, opts DeployOpts) error
	Teardown(ctx context.Context, opts DeployOpts) error
}

type DeployOpts interface {
	ComposeFile(ctx context.Context) (string, error)
	StackName() string
}

type dockerRunner struct {
	Exec func(c *exec.Cmd) error
}

type BuildInfo interface {
	BuildContext() string
	Dockerfile() string
	DockerImage() string
}

func (d dockerRunner) execOutput(c *exec.Cmd) (string, error) {
	buf := new(bytes.Buffer)
	c.Stdout = buf
	err := d.Exec(c)
	return buf.String(), err
}

func (d dockerRunner) SwarmEnabled(ctx context.Context) (bool, error) {
	status, err := d.execOutput(exec.CommandContext(ctx,
		"docker", "info", "--format", "{{ .Swarm.LocalNodeState }}"))
	return status == "active", err
}

func (d dockerRunner) Build(ctx context.Context, context, dockerfile, tag string) error {
	return d.Exec(exec.CommandContext(ctx,
		"docker", "build",
		context,
		"--file", dockerfile,
		"--tag", tag,
	))
}

func (d dockerRunner) Tag(ctx context.Context, src, dst string) error {
	return d.Exec(exec.CommandContext(ctx, "docker", "tag", src, dst))
}

func (d dockerRunner) Push(ctx context.Context, tag string) error {
	return d.Exec(exec.CommandContext(ctx, "docker", "push", tag))
}

func (d dockerRunner) Login(ctx context.Context, hostname, username, password string) error {
	c := exec.CommandContext(ctx,
		"docker", "login",
		hostname,
		"--username", username,
		"--password-stdin",
	)
	c.Stdin = bytes.NewBufferString(password)
	return d.Exec(c)
}

func (d dockerRunner) ComposeConfig(ctx context.Context) (string, error) {
	return d.execOutput(exec.CommandContext(ctx, "docker-compose", "config"))
}
func (d dockerRunner) ComposeUp(ctx context.Context, composeFile, stackName string) error {
	return d.Exec(exec.CommandContext(ctx,
		"docker-compose",
		"--file", composeFile,
		"--project-name", stackName,
		"up",
		"--detach", "--no-build"))
}
func (d dockerRunner) ComposeDown(ctx context.Context, composeFile, stackName string) error {
	return d.Exec(exec.CommandContext(ctx,
		"docker-compose",
		"--file", composeFile,
		"--project-name", stackName,
		"down",
		"--volumes", "--remove-orphans"))
}

func (d dockerRunner) StackDeploy(ctx context.Context, composeFile, stackName string) error {
	return d.Exec(exec.CommandContext(ctx,
		"docker", "stack", "deploy",
		"--compose-file", composeFile,
		"--prune",
		"--with-registry-auth",
		stackName))
}
func (d dockerRunner) StackRemove(ctx context.Context, stackName string) error {
	return d.Exec(exec.CommandContext(ctx, "docker", "stack", "rm", stackName))
}

func writeComposeFile(ctx context.Context, opts DeployOpts) (filename string, err error) {
	contents, err := opts.ComposeFile(ctx)
	if err != nil {
		return "", err
	}
	file, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	_, err = file.WriteString(contents)
	if err != nil {
		return "", err
	}
	return file.Name(), file.Close()
}

func (d dockerRunner) Deploy(ctx context.Context, opts DeployOpts) error {
	stackName := opts.StackName()
	swarmMode, err := d.SwarmEnabled(ctx)
	if err != nil {
		return err
	}

	composeFile, err := writeComposeFile(ctx, opts)
	if err != nil {
		return err
	}
	defer os.Remove(composeFile)

	run := d.ComposeUp
	if swarmMode {
		run = d.StackDeploy
	}
	return run(ctx, composeFile, stackName)
}

func (d dockerRunner) Teardown(ctx context.Context, opts DeployOpts) error {
	stackName := opts.StackName()
	swarmMode, err := d.SwarmEnabled(ctx)
	if err != nil {
		return err
	}

	if swarmMode {
		composeFile, err := writeComposeFile(ctx, opts)
		if err != nil {
			return err
		}
		defer os.Remove(composeFile)
		return d.ComposeDown(ctx, composeFile, stackName)
	} else {
		return d.StackRemove(ctx, stackName)
	}
}

type portainerDeploy struct {
	Client     *portainer.Client
	StackType  string
	EndpointID int64
	GroupID    int64
}

type Errors []error

func (e *Errors) Append(err error) {
	*e = append(*e, err)
}

func (e Errors) Err() error {
	if len(e) > 0 {
		return e
	}
	return nil
}

func (e Errors) Error() string {
	s := ""
	for i, err := range e {
		if i != 0 {
			s += "\n"
		}
		s += err.Error()
	}
	return s
}

func (c *portainerDeploy) tryEndpoints(action func(endpoint portainer.Endpoint) error) error {
	jwt, err := c.Client.Authenticate(nil)
	if err != nil {
		return errors.Wrap(err, "auth")
	}
	c.Client.JWT = jwt

	var endpoints portainer.Endpoints
	if c.EndpointID != 0 {
		var endpoint portainer.Endpoint
		endpoint, err = c.Client.Endpoint(c.EndpointID)
		endpoints = portainer.Endpoints{endpoint}
	} else {
		endpoints, err = c.Client.Endpoints(portainer.EndpointsFilter{GroupID: c.GroupID})
	}
	if err != nil {
		return errors.Wrap(err, "lookup endpoints")
	}

	if len(endpoints) == 0 {
		return errors.New("no applicable endpoints found")
	}

	rand.Shuffle(len(endpoints), endpoints.Swap)

	errors := Errors{}
	for _, endpoint := range endpoints {
		err := action(endpoint)
		if err == nil {
			return nil
		}
		errors.Append(err)
	}
	return errors.Err()
}

func (c *portainerDeploy) findStack(endpointID int64, name string) (*portainer.Stack, error) {
	stacks, err := c.Client.Stacks(portainer.StacksFilter{EndpointID: endpointID})
	if err != nil {
		return nil, errors.Wrap(err, "fetch stacks")
	}
	for _, stack := range stacks {
		if stack.Name == name {
			return &stack, nil
		}
	}
	return nil, nil
}

func (c *portainerDeploy) Deploy(ctx context.Context, opts DeployOpts) error {
	composeFileContents, err := opts.ComposeFile(ctx)
	if err != nil {
		return errors.Wrap(err, "compose file contents")
	}
	name := opts.StackName()

	c.Client.Logger.Printf("Deploying with portainer stack: %v; contents: %v", name, composeFileContents)

	return c.tryEndpoints(func(endpoint portainer.Endpoint) error {
		stack, err := c.findStack(endpoint.ID, name)
		if err != nil {
			return errors.Wrap(err, "find stack")
		}
		if stack == nil {
			s := portainer.Stack{
				EndpointID: endpoint.ID,
				Name:       name,
			}
			_, err = c.Client.NewStack(s, composeFileContents)
			err = errors.Wrap(err, "create stack")
		} else {
			_, err = c.Client.UpdateStack(*stack, composeFileContents, true)
			err = errors.Wrapf(err, "update stack %v", stack.ID)
		}
		return err
	})
}

func (c *portainerDeploy) Teardown(ctx context.Context, opts DeployOpts) error {
	return c.tryEndpoints(func(endpoint portainer.Endpoint) error {
		stack, err := c.findStack(endpoint.ID, opts.StackName())
		if err != nil {
			return err
		}
		if stack != nil {
			err = c.Client.RemoveStack(stack.ID)
		}
		return err
	})
}
