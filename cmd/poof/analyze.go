package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/christopher/poof/internal/analyze"
	"github.com/christopher/poof/internal/config"
	"github.com/christopher/poof/internal/db"
	_ "github.com/christopher/poof/internal/db/postgres"
	"github.com/christopher/poof/internal/ui"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze database schema and suggest columns for masking",
	Run: func(_ *cobra.Command, _ []string) {
		ui.HandleExit(runAnalyze())
	},
}

func runAnalyze() error {
	ctx := context.Background()
	dsn := dbConnStr

	// Try to load config if path exists to support --env
	cfg, err := config.LoadConfig(configPath)
	if err == nil {
		if dsn == "" {
			var dbEnv config.Database
			dbEnv, err = cfg.GetDatabase(envName)
			if err == nil {
				dsn = dbEnv.DSN
			}
		}
	}

	if dsn == "" {
		return ui.WrapError(ui.ErrConfig, config.ErrMissingDSN)
	}

	client, err := db.Connect(ctx, dsn)
	if err != nil {
		return ui.WrapError(ui.ErrConnection, err)
	}
	defer client.Close()

	ui.Info("Analyzing database schema...")
	analyzer := analyze.NewAnalyzer(client)
	suggestions, err := analyzer.Analyze(ctx)
	if err != nil {
		return err
	}

	if len(suggestions) == 0 {
		ui.Success("Analysis complete. No sensitive columns suggested.")
		return nil
	}

	fmt.Println(ui.Bold("\nSuggested Masking Configuration (TOML snippet):"))

	currentTable := ""
	for _, s := range suggestions {
		if s.TableName != currentTable {
			if currentTable != "" {
				fmt.Println()
			}
			fmt.Printf("[[tables]]\nname = %q\npk = \"id\" # Please verify primary key\n", s.TableName)
			currentTable = s.TableName
		}
		fmt.Printf("\n  [[tables.columns]]\n  name = %q\n  [tables.columns.gen]\n  type = %q\n  provider = %q\n", s.ColumnName, s.Generator, s.Provider)
	}

	fmt.Println()
	ui.Success("Analysis complete. Found %d candidates.", len(suggestions))
	ui.Info("Review the suggestions above and copy relevant parts to your poof.toml.")
	return nil
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
