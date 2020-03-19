package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/duckbrain/beluga"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ctx context.Context = context.TODO()
var runner = beluga.New()
var verbose = false

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "beluga",
	Short: "CLI for communicating with the belugad service with detection for CI environments for auto configuration",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		verbose = verbose || runner.DryRun
		if !verbose {
			l := logrus.New()
			l.SetLevel(logrus.FatalLevel)
			runner.Logger = l
		}
		return runner.Env.Validate()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&runner.DryRun, "dry-run", "d", false, "Don't run commands; implies --verbose")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logging of what commands are being run.")
}
