// Package main provides the CLI entry point for poof.
package main

import (
	"context"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/christopher/poof/internal/config"
	"github.com/christopher/poof/internal/db"
	_ "github.com/christopher/poof/internal/db/postgres" // Register Postgres backend
	"github.com/christopher/poof/internal/engine"
	"github.com/christopher/poof/internal/generator"
	"github.com/christopher/poof/internal/ui"
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
		ui.HandleExit(runApply())
	},
}

func runApply() error {
	generator.RegisterAll()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return ui.WrapError(ui.ErrConfig, err)
	}

	ctx := context.Background()
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

	client, err := db.Connect(ctx, dsn)
	if err != nil {
		return ui.WrapError(ui.ErrConnection, err)
	}
	defer client.Close()

	dbName, err := client.GetDatabaseName(ctx)
	if err != nil {
		return ui.WrapError(ui.ErrConnection, err)
	}

	slog.Info("Starting masking process", "database", dbName, "config", configPath, "dry_run", dryRun)

	if !cfg.IsAllowed(dbName) && !force {
		return ui.WrapError(ui.ErrSafety, nil)
	}

	if !yes && !dryRun {
		ui.Info("Running plan before apply...")
		preEngine := engine.NewEngine(client, cfg, workers)
		preEngine.DryRun = true
		_, err = preEngine.Apply(ctx)
		if err != nil {
			return err
		}
		ui.Success("Pre-flight plan completed.")
	}

	if dryRun {
		ui.Info("Applying masking (DRY RUN)...")
	} else {
		ui.Info("Applying masking...")
	}

	eng := engine.NewEngine(client, cfg, workers)
	eng.DryRun = dryRun
	report, err := eng.Apply(ctx)
	if err != nil {
		return err
	}

	totalUpdated := int64(0)
	totalRetried := int64(0)
	totalFailed := int64(0)

	for _, t := range report.Tables {
		totalUpdated += t.Updated
		totalRetried += t.Retried
		totalFailed += t.Failed
	}

	if dryRun {
		ui.Success("Dry run completed successfully. No changes were made.")
	} else {
		ui.Success("Masking completed successfully.")
		slog.Info("Masking process finished.")
	}

	ui.Bold("\nSummary:")
	ui.Info("Updated: %d", totalUpdated)
	ui.Info("Retried: %d (Unique violations resolved)", totalRetried)
	ui.Info("Failed:  %d", totalFailed)

	if totalFailed > 0 {
		return ui.ErrPartial
	}

	return nil
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().BoolVar(&force, "force", false, "Force apply even if database is not in allowlist")
	applyCmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip plan summary and proceed immediately")
	applyCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Execute masking but do not commit changes")
}
