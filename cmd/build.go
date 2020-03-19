package cmd

import (
	"github.com/duckbrain/beluga"
	"github.com/spf13/cobra"
)

var buildOpts beluga.BuildOpts
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds a docker image and pushes it to the registry",
	RunE: func(cmd *cobra.Command, args []string) error {
		must(beluga.New().Build(buildOpts))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().BoolVarP(&buildOpts.Push, "push", "p", true, "If true, push the resulting container to the registry.")
}
