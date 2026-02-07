// Package main provides the CLI entry point for poof.
package main

import (
	"log/slog"
	"os"

	"github.com/christopher/poof/internal/ui"
	"github.com/spf13/cobra"
)

var noColor bool
var configPath string
var dbConnStr string
var workers int

var rootCmd = &cobra.Command{
	Use:   "poof",
	Short: "poof is a PostgreSQL data masking tool",
	Long:  "A declarative, deterministic, and parallel-safe data masking tool for PostgreSQL.",
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		ui.Init(noColor)
		// Initialize structured logging
		handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		slog.SetDefault(slog.New(handler))
	},
}

// Execute starts the CLI application.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("CLI execution failed", "error", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable color output")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "poof.toml", "Path to TOML config file")
	rootCmd.PersistentFlags().StringVarP(&dbConnStr, "db", "d", "", "PostgreSQL connection string (overrides config)")
	rootCmd.PersistentFlags().IntVarP(&workers, "workers", "w", 4, "Number of parallel workers")
}
