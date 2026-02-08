// Package main provides the CLI entry point for poof.
package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/christopher/poof/internal/db"
	_ "github.com/christopher/poof/internal/db/postgres"
	"github.com/christopher/poof/internal/engine"
	"github.com/christopher/poof/internal/generator"
	"github.com/christopher/poof/internal/ui"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Preview masking changes (read-only)",
	Long: `Loads the configuration and database schema to show a detailed plan
of which tables and columns will be masked, including sample transformation previews.
This command never modifies database data.`,
	Run: func(cmd *cobra.Command, _ []string) {
		ui.HandleExit(runPlan(cmd.Context()))
	},
}

func runPlan(ctx context.Context) error {
	generator.RegisterAll()

	cli, err := LoadResources(ctx)
	if err != nil {
		return err
	}
	defer cli.DB.Close()

	// Validate schema
	for _, tableCfg := range cli.Config.Tables {
		var tableCols []db.ColumnInfo
		tableCols, err = cli.DB.GetTableColumns(ctx, tableCfg.Name)
		if err != nil {
			return ui.WrapError(ui.ErrConnection, err)
		}
		if len(tableCols) == 0 {
			return ui.WrapError(ui.ErrConfig, nil)
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

	eng := engine.NewEngine(cli.DB, cli.Config, workers)
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
