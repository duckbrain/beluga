package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate documentation for the CLI",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.MkdirAll(docsPath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		err = doc.GenMarkdownTree(rootCmd, docsPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var docsPath string

func init() {
	docsCmd.Flags().StringVarP(&docsPath, "path", "p", "./docs", "target directory (default: \"./docs\")")
	rootCmd.AddCommand(docsCmd)
}
