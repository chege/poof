package main

import (
	"fmt"
	"os"

	"github.com/christopher/masker/internal/ui"
	"github.com/spf13/cobra"
)

var noColor bool
var configPath string
var dbConnStr string
var workers int

var rootCmd = &cobra.Command{
	Use:   "dbmask",
	Short: "dbmask is a PostgreSQL data masking tool",
	Long:  "A declarative, deterministic, and parallel-safe data masking tool for PostgreSQL.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ui.Init(noColor)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable color output")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "dbmask.hcl", "Path to HCL config file")
	rootCmd.PersistentFlags().StringVarP(&dbConnStr, "db", "d", "", "PostgreSQL connection string")
	rootCmd.PersistentFlags().IntVarP(&workers, "workers", "w", 4, "Number of parallel workers")
}
