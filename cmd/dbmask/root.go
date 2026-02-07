package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dbmask",
	Short: "dbmask is a PostgreSQL data masking tool",
	Long:  "A declarative, deterministic, and parallel-safe data masking tool for PostgreSQL.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Root flags if any
}
