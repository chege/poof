// Package main provides the CLI entry point for poof.
package main

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/christopher/poof/internal/config"
	_ "github.com/christopher/poof/internal/db/postgres"
	"github.com/christopher/poof/internal/engine"
	"github.com/christopher/poof/internal/ui"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check if the environment is ready for masking",
	Run: func(_ *cobra.Command, _ []string) {
		ctx := context.Background()

		// Load config first to get DSN if not provided via flag
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			ui.Error("Config error: %v", err)
			os.Exit(1)
		}

		dsn := dbConnStr
		if dsn == "" {
			dsn = cfg.Database.DSN
		}

		if dsn == "" {
			ui.Error("Database connection string is required (either in config or via --db)")
			os.Exit(1)
		}

		if !engine.CheckReadiness(ctx, configPath, dsn) {
			ui.Error("Doctor found issues. Please fix them before running apply.")
			os.Exit(1)
		}
		ui.Success("Everything looks good! You are ready to mask.")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
