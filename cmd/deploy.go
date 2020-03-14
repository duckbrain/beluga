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
		must(beluga.New().Deploy(ctx))
	},
}

var teardownCmd = &cobra.Command{
	Use:     "teardown",
	Aliases: []string{"beach"},
	Short:   "Teardown an application from a docker daemon",
	Run: func(cmd *cobra.Command, args []string) {
		defer handlePanic()
		must(beluga.New().Teardown(ctx))
	},
}

var composeFileCmd = &cobra.Command{
	Use:   "compose-file",
	Short: "Print out the docker-compose stack that can be deployed",
	Run: func(cmd *cobra.Command, args []string) {
		defer handlePanic()
		contents, err := beluga.New().ComposeFile(ctx)
		must(err)
		fmt.Println(contents)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(teardownCmd)
	rootCmd.AddCommand(composeFileCmd)
}
