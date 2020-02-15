package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/duckbrain/beluga/internal/portainer"
	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var client *portainer.Client
var dsn string
var verbose bool
var endpointsFilter portainer.EndpointsFilter
var stacksFilter portainer.StacksFilter
var stack portainer.Stack
var composeFile string
var prune bool
var id int64
var typeName string

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "portainer",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		c, err := portainer.New(dsn, nil)
		if err != nil {
			return errors.Wrap(err, "parse DSN")
		}
		client = c
		if verbose {
			logger := logrus.New()
			logger.SetLevel(logrus.DebugLevel)
			client.Logger = logger
		}

		client.Logger.Println("logging in")
		jwt, err := client.Authenticate(nil)
		if err != nil {
			return errors.Wrap(err, "authentication")
		}
		client.JWT = jwt

		stack.Env = portainer.Env(envy.Map())

		return nil
	},
}

var endpointsCmd = &cobra.Command{
	Use:     "endpoints",
	Aliases: []string{"endpoint", "e"},
	Short:   "List/view endpoints",
	RunE: func(cmd *cobra.Command, args []string) error {
		var endpoints portainer.Endpoints
		var err error
		if len(args) > 0 {
			for i, arg := range args {
				id, err := strconv.ParseInt(arg, 10, 64)
				if err != nil {
					return errors.Wrapf(err, "parse endpoint ID %v", i+1)
				}
				e, err := client.Endpoint(id)
				if err != nil {
					return err
				}
				endpoints = append(endpoints, e)
			}
		} else {
			endpoints, err = client.Endpoints(endpointsFilter)
		}
		if err != nil {
			return err
		}
		fmt.Println(endpoints)
		return nil
	},
}

var stacksCmd = &cobra.Command{
	Use:     "stacks",
	Aliases: []string{"stack", "s"},
	Short:   "List/view stacks",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return errors.New("Lookup by ID not yet supported")
		}

		stacks, err := client.Stacks(stacksFilter)
		if err != nil {
			return err
		}
		fmt.Println(stacks)
		return nil
	},
}

func preIDArg(cmd *cobra.Command, args []string) (err error) {
	if len(args) == 0 {
		return errors.New("must provide an ID")
	}
	if len(args) > 1 {
		return errors.New("can only provide one ID")
	}
	id, err = strconv.ParseInt(args[0], 10, 64)
	return
}

var stackNewCmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "n", "c"},
	Short:   "Create a new stack",
	RunE: func(cmd *cobra.Command, args []string) error {
		composeFileContents, err := ComposeFileContents()
		if err != nil {
			return errors.Wrap(err, "read compose file")
		}
		stack.Type, err = portainer.ParseStackType(typeName)
		if err != nil {
			return err
		}
		s, err := client.NewStack(stack, composeFileContents)
		if err != nil {
			return err
		}
		fmt.Println(s)
		return nil
	},
}
var stackUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"create", "n", "c"},
	Short:   "Create a new stack",
	PreRunE: preIDArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		stack.ID = id
		composeFileContents, err := ComposeFileContents()
		if err != nil {
			return errors.Wrap(err, "read compose file")
		}
		s, err := client.UpdateStack(stack, composeFileContents, prune)
		if err != nil {
			return err
		}
		fmt.Println(s)
		return nil
	},
}
var stackRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"delete", "destroy", "drop", "r", "d"},
	Short:   "Remove a new stack",
	PreRunE: preIDArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.RemoveStack(id)
	},
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func ComposeFileContents() (string, error) {
	var data []byte
	var err error
	if composeFile == "-" {
		client.Logger.Println("reading compose from stdin")
		data, err = ioutil.ReadAll(os.Stdin)
	} else {
		client.Logger.Printf("reading compose file %v", composeFile)
		data, err = ioutil.ReadFile(composeFile)
	}
	return string(data), err
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&dsn, "dsn", "d", envy.Get("PORTAINER_DSN", ""), "DNS to connect to portainer with")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print out debugging logging")

	endpointsCmd.Flags().Int64VarP(&endpointsFilter.GroupID, "group-id", "g", 0, "Filter by group ID")
	rootCmd.AddCommand(endpointsCmd)

	stacksCmd.Flags().Int64VarP(&stacksFilter.EndpointID, "endpoint-id", "e", 0, "Filter by endpoint ID")
	rootCmd.AddCommand(stacksCmd)

	stackNewCmd.Flags().Int64VarP(&stack.EndpointID, "endpoint-id", "e", 0, "Endpoint ID to deploy the stack onto")
	stackNewCmd.Flags().StringVarP(&stack.Name, "name", "n", "", "Name of the stack")
	stackNewCmd.Flags().StringVarP(&typeName, "type", "t", "swarm", "Type of stack: swarm | compose")
	stackNewCmd.Flags().StringVarP(&composeFile, "compose-file", "f", "docker-compose.yaml", `Filename for the compose file to deploy; if "-", will use stdin`)
	must(stackNewCmd.MarkFlagRequired("name"))
	must(stackNewCmd.MarkFlagRequired("endpoint-id"))
	stacksCmd.AddCommand(stackNewCmd)

	stackUpdateCmd.Flags().Int64VarP(&stack.EndpointID, "endpoint-id", "e", 0, "Endpoint ID to deploy the stack onto")
	stackUpdateCmd.Flags().BoolVarP(&prune, "prune", "p", true, "Prune the containers when updating the stack")
	stacksCmd.AddCommand(stackUpdateCmd)

	stacksCmd.AddCommand(stackRemoveCmd)
}
