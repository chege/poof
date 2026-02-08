// Package main provides the CLI entry point for poof.
package main

import (
	"context"

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
		ui.HandleExit(runPlan())
	},
}

func runPlan() error {
	generator.RegisterAll()
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

	ctx := context.Background()
	client, err := db.Connect(ctx, dsn)
	if err != nil {
		return ui.WrapError(ui.ErrConnection, err)
	}
	defer client.Close()

	// Validate schema
	for _, tableCfg := range cfg.Tables {
		var tableCols []db.ColumnInfo
		tableCols, err = client.GetTableColumns(ctx, tableCfg.Name)
		if err != nil {
			return ui.WrapError(ui.ErrConnection, err)
		}
		if len(tableCols) == 0 {
			return ui.WrapError(ui.ErrConfig, nil) // Could use a more specific error
		}

		colMap := make(map[string]bool)
		for _, c := range tableCols {
			colMap[c.Name] = true
		}

		for _, c := range tableCfg.Columns {
			if !colMap[c.Name] {
				return ui.WrapError(ui.ErrConfig, nil)
			}
		}
	}

	eng := engine.NewEngine(client, cfg, workers)
	eng.DryRun = true

	ui.Bold("Masking Plan:")
	report, err := eng.Apply(ctx)
	if err != nil {
		return err
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
	return nil
}

func init() {
	rootCmd.AddCommand(planCmd)
}
