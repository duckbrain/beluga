package cmd

import (
	"fmt"

	"github.com/duckbrain/beluga"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:     "deploy",
	Aliases: []string{"breach"},
	Short:   "Deploy an application to a docker daemon",
	Run: func(cmd *cobra.Command, args []string) {
		defer handlePanic()
		env := beluga.Env()
		must(env.Deployer().Deploy(env))
	},
}

var teardownCmd = &cobra.Command{
	Use:     "teardown",
	Aliases: []string{"beach"},
	Short:   "Teardown an application from a docker daemon",
	Run: func(cmd *cobra.Command, args []string) {
		defer handlePanic()
		env := beluga.Env()
		must(env.Deployer().Teardown(env))
	},
}

var stackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Print out the docker-compose stack that can be deployed",
	Run: func(cmd *cobra.Command, args []string) {
		defer handlePanic()
		env := beluga.Env()
		contents, err := env.ComposeFileContents()
		must(err)
		fmt.Println(contents)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(teardownCmd)
	rootCmd.AddCommand(stackCmd)
}
