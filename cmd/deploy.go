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
		must(lib.Env().Deployer().Deploy())
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
