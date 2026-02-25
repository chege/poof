// Package main provides the CLI entry point for poof.
package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/chege/poof/internal/engine"
	"github.com/chege/poof/internal/generator"
	"github.com/chege/poof/internal/ui"
)

var (
	force      bool
	yes        bool
	dryRun     bool
	reportPath string
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Execute data masking on the database",
	Long: `Applies the configured masking rules to the selected database environment.
Supports --dry-run to simulate the transformation without committing changes.`,
	Run: func(cmd *cobra.Command, _ []string) {
		ui.HandleExit(runApply(cmd.Context()))
	},
}

func runApply(ctx context.Context) error {
	generator.RegisterAll()

	cli, err := LoadResources(ctx)
	if err != nil {
		return err
	}
	defer cli.DB.Close()

	dbName, err := cli.DB.GetDatabaseName(ctx)
	if err != nil {
		return ui.WrapError(ui.ErrConnection, err)
	}

	slog.Info("Starting masking process", "database", dbName, "config", configPath, "dry_run", dryRun)

	if !cli.Config.IsAllowed(dbName) && !force {
		return ui.WrapError(ui.ErrSafety, nil)
	}

	if !yes && !dryRun {
		ui.Info("Running plan before apply...")
		preEngine := engine.NewEngine(cli.DB, cli.Config, workers)
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

	eng := engine.NewEngine(cli.DB, cli.Config, workers)
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
	ui.Info("Duration: %v", report.Duration.Round(time.Millisecond))

	if reportPath != "" {
		f, err := os.Create(filepath.Clean(reportPath))
		if err != nil {
			ui.Error("Failed to create report file: %v", err)
		} else {
			encoder := json.NewEncoder(f)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(report); err != nil {
				ui.Error("Failed to write report: %v", err)
			} else {
				ui.Success("Execution report saved to %s", reportPath)
			}
			_ = f.Close()
		}
	}

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
	applyCmd.Flags().StringVar(&reportPath, "report", "", "Path to save the masking execution report (JSON)")
}
