package cmd

import (
	"fmt"

	"github.com/duckbrain/beluga"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var logger = logrus.New()

var envAllKeys bool
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Print out the environment variables that beluga computes.",
	Run: func(cmd *cobra.Command, args []string) {
		values, err := beluga.Env().Format(beluga.BashFormat, envAllKeys)
		if err != nil {
			if values == nil {
				logger.Error(err)
				return
			} else {
				logger.Warn(err)
			}
		}
		for _, line := range values {
			fmt.Println(line)
		}
	},
}

func init() {
	envCmd.Flags().BoolVarP(&envAllKeys, "all", "a", false, "Output all environment variables, instead of just those known to beluga.")
	rootCmd.AddCommand(envCmd)
}
