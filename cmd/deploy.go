package cmd

import (
	"github.com/duckbrain/beluga/internal/docker"
	"github.com/duckbrain/beluga/internal/lib"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:     "deploy",
	Aliases: []string{"beach"},
	Short:   "Deploy an application to a docker daemon",
	Run: func(cmd *cobra.Command, args []string) {
		defer handlePanic()
		env := lib.Env()
		compose := docker.New(env.DeployDockerHost()).Compose(env)
		must(compose.Deploy())
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
