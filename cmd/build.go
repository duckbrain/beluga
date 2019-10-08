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
		builtImage := ""
		images := strings.Fields(e.Image())
		d := docker.New("") // TODO: Have a way to specify for build

		for _, image := range images {
			if builtImage == "" {
				must(d.Build(e.DockerContext(), e.Dockerfile(), image))
				builtImage = image
			} else {
				must(d.Tag(builtImage, image))
			}
		}
		must(d.Push(images...))
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
