package main

import (
	"fmt"
	"os"

	"github.com/duckbrain/beluga/portainer"
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

		return nil
	},
}

var endpointsCmd = &cobra.Command{
	Use:     "endpoints",
	Aliases: []string{"endpoint", "e"},
	Short:   "List/view endpoints",
	RunE: func(cmd *cobra.Command, args []string) error {
		endpoints, err := client.Endpoints(endpointsFilter)
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
		stacks, err := client.Stacks(stacksFilter)
		if err != nil {
			return err
		}
		fmt.Println(stacks)
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&dsn, "dsn", "d", envy.Get("PORTAINER_DSN", ""), "DNS to connect to portainer with")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print out debugging logging")

	endpointsCmd.Flags().Int64VarP(&endpointsFilter.GroupID, "group-id", "g", 0, "Filter by group ID")
	rootCmd.AddCommand(endpointsCmd)

	stacksCmd.Flags().Int64VarP(&stacksFilter.EndpointID, "endpoint-id", "e", 0, "Filter by endpoint ID")
	rootCmd.AddCommand(stacksCmd)
}
