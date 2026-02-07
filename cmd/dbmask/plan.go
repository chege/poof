package main

import (
	"context"
	"os"

	"github.com/christopher/masker/internal/config"
	"github.com/christopher/masker/internal/db"
	"github.com/christopher/masker/internal/generator"
	"github.com/christopher/masker/internal/masker"
	"github.com/christopher/masker/internal/ui"
	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Show a summary of masking changes",
	Run: func(cmd *cobra.Command, args []string) {
		generator.RegisterAll()
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			ui.Error("Config error: %v", err)
			os.Exit(1)
		}

		ctx := context.Background()
		if dbConnStr == "" {
			ui.Error("Database connection string is required (--db)")
			os.Exit(1)
		}

		client, err := db.NewClient(ctx, dbConnStr)
		if err != nil {
			ui.Error("DB error: %v", err)
			os.Exit(1)
		}
		defer client.Close()

		engine := masker.NewEngine(client, cfg, workers)
		engine.DryRun = true

		ui.Bold("Masking Plan:")
		diffs, err := engine.Apply(ctx)
		if err != nil {
			ui.Error("Plan failed: %v", err)
			os.Exit(1)
		}

		currentTable := ""
		for _, d := range diffs {
			if d.TableName != currentTable {
				ui.Info("Table: %s", ui.Bold(d.TableName))
				currentTable = d.TableName
			}
			ui.Info("  [%v] %s: %v -> %v", d.PKValue, d.ColumnName, d.OldValue, d.NewValue)
		}

		ui.Success("Plan generated for %d changes.", len(diffs))
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
