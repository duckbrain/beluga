package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:     "deploy",
	Aliases: []string{"breach", "publish"},
	Short:   "Deploy an application to a docker daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.Deploy(ctx)
	},
}

var teardownCmd = &cobra.Command{
	Use:     "teardown",
	Aliases: []string{"beach"},
	Short:   "Teardown an application from a docker daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.Teardown(ctx)
	},
}

var composeFileCmd = &cobra.Command{
	Use:   "compose-file",
	Short: "Print out the docker-compose stack that can be deployed",
	RunE: func(cmd *cobra.Command, args []string) error {
		contents, err := runner.ComposeFile(ctx)
		if err != nil {
			return err
		}
		fmt.Println(contents)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(teardownCmd)
	rootCmd.AddCommand(composeFileCmd)
}
