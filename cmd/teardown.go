package cmd

import (
	"github.com/duckbrain/beluga/internal/lib"
	"github.com/spf13/cobra"
)

var teardownCmd = &cobra.Command{
	Use:     "teardown",
	Aliases: []string{"beach"},
	Short:   "Teardown an application from a docker daemon",
	Run: func(cmd *cobra.Command, args []string) {
		defer handlePanic()
		must(lib.Env().Deployer().Teardown())
	},
}

func init() {
	rootCmd.AddCommand(teardownCmd)
}
