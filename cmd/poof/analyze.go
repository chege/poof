package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/christopher/poof/internal/analyze"
	_ "github.com/christopher/poof/internal/db/postgres"
	"github.com/christopher/poof/internal/ui"
)

var outputJSON bool

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze database schema and suggest columns for masking",
	Run: func(cmd *cobra.Command, _ []string) {
		ui.HandleExit(runAnalyze(cmd.Context()))
	},
}

func runAnalyze(ctx context.Context) error {
	cli, err := LoadResources(ctx)
	if err != nil {
		return err
	}
	defer cli.DB.Close()

	if !outputJSON {
		ui.Info("Analyzing database schema...")
	}

	analyzer := analyze.NewAnalyzer(cli.DB)
	suggestions, err := analyzer.Analyze(ctx)
	if err != nil {
		return err
	}

	if outputJSON {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(suggestions)
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
	analyzeCmd.Flags().BoolVar(&outputJSON, "json", false, "Output results in JSON format")
}
