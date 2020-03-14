package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var ctx context.Context = context.TODO()

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func handlePanic() {
	msg := recover()
	if err, ok := msg.(error); ok {
		fmt.Printf("beluga: %v\n", err)
		os.Exit(1)
	} else if msg != nil {
		fmt.Printf("beluga: unknown error: %v\n", msg)
		os.Exit(2)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "beluga",
	Short: "CLI for communicating with the belugad service with detection for CI environments for auto configuration",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.beluga.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
