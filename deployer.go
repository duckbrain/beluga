package beluga

import (
	"context"
	"io/ioutil"
	"math/rand"
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

func (r Runner) deployer() Deployer {
	host := r.Env["DOCKER_HOST"]

	switch host[:strings.Index(host, ":")] {
	case "portainer", "portainer-insecure":
		deployer := &PortainerDeploy{}
		client, err := portainer.New(host, deployer)
		if err != nil {
			panic(err)
		}
		client.Logger = r.Logger
		deployer.Client = client
		return deployer
	case "ssh":
		panic("SSH not implemented")
	default:
		return dockerRunner{Exec: r.Exec}
	}
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

func (d dockerRunner) Build(context, dockerfile, tag string) error {
	return d.Exec(exec.Command(
		"docker", "build",
		context,
		"-f", dockerfile,
		"-t", tag,
	))
}

func (d dockerRunner) Tag(src, dst string) error {
	return d.Exec(exec.Command("docker", "tag", src, dst))
}

func (d dockerRunner) Push(tag string) error {
	return d.Exec(exec.Command("docker", "push", tag))
}

func (d dockerRunner) Login(hostname, username, password string) error {
	return d.Exec(exec.Command(
		"docker", "login",
		hostname,
		"-u", username,
		"-p", password,
	))
}

func writeComposeFile(ctx context.Context, opts DeployOpts) (filename string, err error) {
	contents, err := opts.ComposeFile(ctx)
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
	composeFile, err := writeComposeFile(ctx, opts)
	if err != nil {
		return err
	}
	defer os.Remove(composeFile)

	var cmd *exec.Cmd
	switch opts.DeployMode() {
	case ComposeMode:
		cmd = exec.Command("docker-compose",
			"--file", composeFile,
			"--project-name", stackName,
			"up",
			"--detach", "--no-build")
	case SwarmMode:
		cmd = exec.Command("docker", "stack", "deploy",
			"--compose-file", composeFile,
			"--prune",
			"--with-registry-auth",
			stackName)
	}
	return d.run(cmd)
}

func (d dockerRunner) Teardown(ctx context.Context, opts DeployOpts) error {
	stackName := opts.StackName()
	var cmd *exec.Cmd
	switch opts.DeployMode() {
	case ComposeMode:
		composeFile, err := writeComposeFile(ctx, opts)
		if err != nil {
			return err
		}
		defer os.Remove(composeFile)
		cmd = exec.Command("docker-compose",
			"--file", composeFile,
			"--project-name", stackName,
			"down",
			"--volumes", "--remove-orphans")
	case SwarmMode:
		cmd = exec.Command("docker", "stack", "rm", stackName)
	}
	return d.run(cmd)
}

type PortainerDeploy struct {
	Client     *portainer.Client
	StackType  string
	EndpointID int64
	GroupID    int64
}

type Errors []error

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

func (c *PortainerDeploy) tryEndpoints(action func(endpoint portainer.Endpoint) error) error {
	jwt, err := c.Client.Authenticate(nil)
	if err != nil {
		errors.Wrap(err, "auth")
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
		errors = append(errors, err)
	}
	return errors
}

func (c *PortainerDeploy) findStack(endpointID int64, name string) (*portainer.Stack, error) {
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

func (c *PortainerDeploy) Deploy(ctx context.Context, opts DeployOpts) error {
	composeFileContents, err := opts.ComposeFile(ctx)
	if err != nil {
		return errors.Wrap(err, "compose file contents")
	}
	stackType := portainer.StackType(opts.DeployMode())
	name := opts.StackName()
	c.Client.Logger.Printf("Deploying with portainer %v in %v\n%v", name, stackType, composeFileContents)

	return c.tryEndpoints(func(endpoint portainer.Endpoint) error {
		stack, err := c.findStack(endpoint.ID, name)
		if err != nil {
			return errors.Wrap(err, "find stack")
		}
		if stack == nil {
			s := portainer.Stack{
				EndpointID: endpoint.ID,
				Name:       name,
				Type:       stackType,
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

func (c *PortainerDeploy) Teardown(ctx context.Context, opts DeployOpts) error {
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
