package cmd

import (
	"strings"

	"github.com/duckbrain/beluga/internal/docker"
	"github.com/duckbrain/beluga/internal/lib"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds a docker image and pushes it to the registry",
	Run: func(cmd *cobra.Command, args []string) {
		defer handlePanic()
		e := lib.Env()

		if e.RegistryUsername() != "" {
			must(docker.Login(e.Registry(), e.RegistryUsername(), e.RegistryPassword()))
		}

		builtImage := ""
		images := strings.Fields(e.Image())
		for _, image := range images {
			if builtImage == "" {
				must(docker.Build(e.DockerContext(), e.Dockerfile(), image))
				builtImage = image
			} else {
				must(docker.Tag(builtImage, image))
			}
		}
		must(docker.Push(images...))
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
