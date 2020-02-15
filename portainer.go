package beluga

import (
	"math/rand"

	"github.com/duckbrain/beluga/internal/portainer"
	"github.com/pkg/errors"
)

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
	var endpoints portainer.Endpoints
	var err error
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

func (c *PortainerDeploy) Deploy(opts DeployOpts) error {
	composeFileContents, err := opts.ComposeFileContents()
	if err != nil {
		return errors.Wrap(err, "compose file contents")
	}
	stackType := portainer.StackType(opts.DeployMode())

	return c.tryEndpoints(func(endpoint portainer.Endpoint) error {
		stack, err := c.findStack(endpoint.ID, opts.StackName())
		if err != nil {
			return err
		}
		if stack == nil {
			s := portainer.Stack{
				EndpointID: endpoint.ID,
				Name:       opts.StackName(),
				Type:       stackType,
			}
			_, err = c.Client.NewStack(s, composeFileContents)
		} else {
			_, err = c.Client.UpdateStack(*stack, composeFileContents, true)
		}
		return err
	})
}

func (c *PortainerDeploy) Teardown(opts DeployOpts) error {
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
