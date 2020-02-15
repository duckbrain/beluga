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

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "portainer",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		c, err := portainer.New(dsn)
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
		endpoints, err := client.ListEndpoints()
		if err != nil {
			return err
		}
		fmt.Println(endpoints)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(endpointsCmd)
	rootCmd.PersistentFlags().StringVarP(&dsn, "dsn", "d", envy.Get("PORTAINER_DSN", ""), "DNS to connect to portainer with")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print out debugging logging")
}
