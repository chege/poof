// Package main provides the CLI entry point for poof.
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/christopher/poof/internal/config"
	"github.com/christopher/poof/internal/db"
	"github.com/christopher/poof/internal/ui"
)

var noColor bool
var configPath string
var dbConnStr string
var workers int
var envName string

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
		// Cobra already prints the error, but we want to ensure it's handled via our UI
		os.Exit(1)
	}
}

// CLIContext provides shared resources for CLI commands.
type CLIContext struct {
	Config *config.Config
	DB     db.DB
}

// LoadResources loads the configuration and connects to the database based on flags.
func LoadResources(ctx context.Context) (*CLIContext, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, ui.WrapError(ui.ErrConfig, err)
	}

	dsn := dbConnStr
	if dsn == "" {
		var dbEnv config.Database
		dbEnv, err = cfg.GetDatabase(envName)
		if err != nil {
			return nil, ui.WrapError(ui.ErrConfig, err)
		}
		dsn = dbEnv.DSN
	}

	if dsn == "" {
		return nil, ui.WrapError(ui.ErrConfig, config.ErrMissingDSN)
	}

	client, err := db.Connect(ctx, dsn)
	if err != nil {
		return nil, ui.WrapError(ui.ErrConnection, err)
	}

	return &CLIContext{
		Config: cfg,
		DB:     client,
	}, nil
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable color output")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "poof.toml", "Path to TOML config file")
	rootCmd.PersistentFlags().StringVarP(&dbConnStr, "db", "d", "", "PostgreSQL connection string (overrides config)")
	rootCmd.PersistentFlags().StringVarP(&envName, "env", "e", "", "Select database environment from config")
	rootCmd.PersistentFlags().IntVarP(&workers, "workers", "w", 4, "Number of parallel workers")
}
