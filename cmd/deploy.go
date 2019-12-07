package cmd

import (
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
		must(env.Deployer().Deploy(env.DeployOpts()))
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
