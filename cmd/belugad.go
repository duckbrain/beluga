package cmd

import (
	"net/http"

	"github.com/duckbrain/beluga/internal/lib"
)

var dsn string

func belugaClient() *lib.Client {
	return &lib.Client{
		DSN:    dsn,
		Client: http.DefaultClient,
	}
}

func init() {
	rootCmd.Flags().StringVarP(&dsn, "dsn", "d", "", "target DSN")
}
