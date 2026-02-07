// Package main provides the CLI entry point for poof.
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/christopher/poof/internal/config"
	"github.com/christopher/poof/internal/db"
	_ "github.com/christopher/poof/internal/db/postgres" // Register Postgres backend
	"github.com/christopher/poof/internal/engine"
	"github.com/christopher/poof/internal/generator"
	"github.com/christopher/poof/internal/ui"
	"github.com/spf13/cobra"
)

var (
	force  bool
	yes    bool
	dryRun bool
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply data masking rules to the database",
	Run: func(_ *cobra.Command, _ []string) {
		generator.RegisterAll()

		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			ui.Error("Config error: %v", err)
			os.Exit(1)
		}

		ctx := context.Background()
		if dbConnStr == "" {
			dbConnStr = cfg.Database.DSN
		}

		if dbConnStr == "" {
			ui.Error("Database connection string is required (either in config or via --db)")
			os.Exit(1)
		}

		client, err := db.Connect(ctx, dbConnStr)
		if err != nil {
			ui.Error("DB connection error: %v", err)
			os.Exit(1)
		}
		defer client.Close()

		dbName, err := client.GetDatabaseName(ctx)
		if err != nil {
			ui.Error("DB name error: %v", err)
			os.Exit(1)
		}

		slog.Info("Starting masking process", "database", dbName, "config", configPath, "dry_run", dryRun)

		if !cfg.IsAllowed(dbName) && !force {
			ui.Error("Database %q is not in the allowed_db_names list and --force was not provided.", dbName)
			os.Exit(1)
		}

		if !yes && !dryRun {
			ui.Info("Running plan before apply...")
			preEngine := engine.NewEngine(client, cfg, workers)
			preEngine.DryRun = true
			_, err = preEngine.Apply(ctx)
			if err != nil {
				ui.Error("Pre-apply plan failed: %v", err)
				os.Exit(1)
			}
			ui.Success("Pre-flight plan completed.")
		}

		if dryRun {
			ui.Info("Applying masking (DRY RUN)...")
		} else {
			ui.Info("Applying masking...")
		}

		engine := engine.NewEngine(client, cfg, workers)
		engine.DryRun = dryRun
		_, err = engine.Apply(ctx)
		if err != nil {
			ui.Error("Apply failed: %v", err)
			os.Exit(1)
		}

		if dryRun {
			ui.Success("Dry run completed successfully. No changes were made.")
		} else {
			ui.Success("Masking completed successfully.")
			slog.Info("Masking process finished.")
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().BoolVar(&force, "force", false, "Force apply even if database is not in allowlist")
	applyCmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip plan summary and proceed immediately")
	applyCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Execute masking but do not commit changes")
}
