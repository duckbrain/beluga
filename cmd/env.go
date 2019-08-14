package cmd

import (
	"fmt"

	"github.com/duckbrain/beluga/internal/lib"
	"github.com/spf13/cobra"
)

var envAllKeys bool
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Print out the environment variables that beluga computes.",
	Run: func(cmd *cobra.Command, args []string) {
		e := lib.Environment
		var keys []string
		if envAllKeys {
			keys = e.SortedKeys()
		} else {
			keys = e.KnownKeys()
		}
		for _, key := range keys {
			fmt.Printf("export %v='%v'\n", key, e[key])
		}
	},
}

func init() {
	envCmd.Flags().BoolVarP(&envAllKeys, "all", "a", false, "Output all environment variables, instead of just those known to beluga.")
	rootCmd.AddCommand(envCmd)
}
