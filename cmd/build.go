package cmd

import (
	"strings"

	"github.com/duckbrain/beluga"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds a docker image and pushes it to the registry",
	Run: func(cmd *cobra.Command, args []string) {
		defer handlePanic()
		e := beluga.Env()
		d := beluga.Docker("") // TODO: Have a way to specify for build

		if e.RegistryUsername() != "" {
			must(d.Login(e.Registry(), e.RegistryUsername(), e.RegistryPassword()))
		}

		builtImage := ""
		images := strings.Fields(e.Image())

		for _, image := range images {
			if builtImage == "" {
				must(d.Build(e.DockerContext(), e.Dockerfile(), image))
				builtImage = image
			} else {
				must(d.Tag(builtImage, image))
			}
		}
		for _, image := range images {
			must(d.Push(image))
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
