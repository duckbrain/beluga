package cmd

import (
	"github.com/duckbrain/beluga/internal/lib"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:     "deploy",
	Aliases: []string{"breach"},
	Short:   "Deploy an application to a docker daemon",
	Run: func(cmd *cobra.Command, args []string) {
		defer handlePanic()
		env := lib.Env()
		must(env.Deployer().Deploy(env.DockerComposeFile(), map[string]string(env)))
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
