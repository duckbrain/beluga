package cmd

import (
	"fmt"

	"github.com/duckbrain/beluga/internal/lib"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Print out the environment variables that beluga computes.",
	Run: func(cmd *cobra.Command, args []string) {
		env := lib.ParseEnv()
		for _, key := range env.KnownKeys() {
			val := env[key]
			fmt.Printf("export %v='%v'\n", key, val)
		}
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
}
