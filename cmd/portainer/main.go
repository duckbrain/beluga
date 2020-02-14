package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "portainer",
}

var endpointsCmd = &cobra.Command{
	Use:     "endpoints",
	Aliases: []string{"endpoint", "e"},
	Short:   "List/view endpoints",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("endpoints here")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(endpointsCmd)
}
