package cmd

import (
	"strings"

	"github.com/duckbrain/beluga"
	"github.com/spf13/cobra"
)

var buildOpts struct {
	// Cache bool // TODO
	Push bool
}
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds a docker image and pushes it to the registry",
	RunE: func(cmd *cobra.Command, args []string) error {
		e := beluga.Env()
		d := beluga.Docker("") // TODO: Have a way to specify for build

		if e.RegistryUsername() != "" {
			err := d.Login(e.Registry(), e.RegistryUsername(), e.RegistryPassword())
			if err != nil {
				return err
			}
		}

		builtImage := ""
		images := strings.Fields(e.Image())

		for _, image := range images {
			if builtImage == "" {
				err := d.Build(e.DockerContext(), e.Dockerfile(), image)
				if err != nil {
					return err
				}
				builtImage = image
			} else {
				err := d.Tag(builtImage, image)
				if err != nil {
					return err
				}
			}
		}

		if buildOpts.Push {
			for _, image := range images {
				err := d.Push(image)
				if err != nil {
					return err
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().BoolVarP(&buildOpts.Push, "push", "p", true, "If true, push the resulting container to the registry.")
}
