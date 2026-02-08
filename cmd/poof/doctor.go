// Package main provides the CLI entry point for poof.
package main

import (
	"context"

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
		ui.HandleExit(runDoctor())
	},
}

func runDoctor() error {
	ctx := context.Background()

	// Load config first to get DSN if not provided via flag
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return ui.WrapError(ui.ErrConfig, err)
	}

	dsn := dbConnStr
	if dsn == "" {
		var dbEnv config.Database
		dbEnv, err = cfg.GetDatabase(envName)
		if err != nil {
			return ui.WrapError(ui.ErrConfig, err)
		}
		dsn = dbEnv.DSN
	}

	if dsn == "" {
		return ui.WrapError(ui.ErrConfig, config.ErrMissingDSN)
	}

	if !engine.CheckReadiness(ctx, configPath, dsn) {
		return ui.WrapError(ui.ErrConnection, nil)
	}
	ui.Success("Everything looks good! You are ready to mask.")
	return nil
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
