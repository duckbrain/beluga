package cmd

import (
	"errors"
	"fmt"

	"github.com/duckbrain/beluga"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var logger = logrus.New()

var envOpts struct {
	All    bool
	Format string
}
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Print out the environment variables that beluga computes.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var format beluga.EnvFormat
		switch envOpts.Format {
		case "env":
			format = beluga.EnvFileFormat
		case "bash":
			format = beluga.BashFormat
		default:
			return errors.New("unknown format")
		}

		values, err := beluga.Env().Format(format, envOpts.All)
		if err != nil {
			if values == nil {
				return err
			} else {
				logger.Warn(err)
			}
		}
		for _, line := range values {
			fmt.Println(line)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
	envCmd.Flags().BoolVarP(&envOpts.All, "all", "a", false, "Output all environment variables, instead of just those known to beluga.")
	envCmd.Flags().StringVarP(&envOpts.Format, "format", "f", "bash", "Format to output environment values, bash|env")
}
