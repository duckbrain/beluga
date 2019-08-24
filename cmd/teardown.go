package cmd

import (
	"github.com/duckbrain/beluga/internal/docker"
	"github.com/duckbrain/beluga/internal/lib"
	"github.com/spf13/cobra"
)

var teardownCmd = &cobra.Command{
	Use:     "teardown",
	Aliases: []string{"beach"},
	Short:   "Teardown an application from a docker daemon",
	Run: func(cmd *cobra.Command, args []string) {
		defer handlePanic()
		compose := docker.Compose(lib.Env())
		must(compose.Teardown())
	},
}

func init() {
	rootCmd.AddCommand(teardownCmd)
}
