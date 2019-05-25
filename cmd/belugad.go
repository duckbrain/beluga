package cmd

import (
	"net/http"

	"github.com/duckbrain/beluga/beluga"
)

var dsn string

func belugaClient() *beluga.Client {
	return &beluga.Client{
		DSN:    dsn,
		Client: http.DefaultClient,
	}
}

func init() {
	rootCmd.Flags().StringVarP(&dsn, "dsn", "d", "", "target DSN")
}
