// Package main provides the CLI entry point for poof.
package main

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/christopher/poof/internal/config"
	"github.com/christopher/poof/internal/db"
	_ "github.com/christopher/poof/internal/db/postgres"
	"github.com/christopher/poof/internal/engine"
	"github.com/christopher/poof/internal/generator"
	"github.com/christopher/poof/internal/ui"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Show a summary of masking changes",
	Run: func(_ *cobra.Command, _ []string) {
		generator.RegisterAll()
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

		ctx := context.Background()
		client, err := db.Connect(ctx, dsn)
		if err != nil {
			ui.Error("DB error: %v", err)
			os.Exit(1)
		}
		defer client.Close()

		engine := engine.NewEngine(client, cfg, workers)
		engine.DryRun = true

		ui.Bold("Masking Plan:")
		report, err := engine.Apply(ctx)
		if err != nil {
			ui.Error("Plan failed: %v", err)
			os.Exit(1)
		}

		for _, t := range report.Tables {
			ui.Info("Table: %s (~%d rows)", ui.Bold(t.Name), t.Estimates)
			for _, col := range t.Columns {
				ui.Info("  - Column: %s", col)
			}
		}

		ui.Bold("\nSample Changes (first 5 rows):")
		currentTable := ""
		for _, d := range report.Diffs {
			if d.TableName != currentTable {
				ui.Info("Table: %s", ui.Bold(d.TableName))
				currentTable = d.TableName
			}
			ui.Info("  [%v] %s: %v -> %v", d.PKValue, d.ColumnName, d.OldValue, d.NewValue)
		}

		ui.Success("Plan generated successfully.")
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
