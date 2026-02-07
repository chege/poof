package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/christopher/poof/internal/analyze"
	"github.com/christopher/poof/internal/db"
	_ "github.com/christopher/poof/internal/db/postgres"
	"github.com/christopher/poof/internal/ui"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze database schema and suggest columns for masking",
	Run: func(_ *cobra.Command, _ []string) {
		if dbConnStr == "" {
			ui.Error("Database connection string is required via --db")
			os.Exit(1)
		}

		ctx := context.Background()
		client, err := db.Connect(ctx, dbConnStr)
		if err != nil {
			ui.Error("DB error: %v", err)
			os.Exit(1)
		}
		defer client.Close()

		ui.Info("Analyzing database schema...")
		analyzer := analyze.NewAnalyzer(client)
		suggestions, err := analyzer.Analyze(ctx)
		if err != nil {
			ui.Error("Analysis failed: %v", err)
			os.Exit(1)
		}

		if len(suggestions) == 0 {
			ui.Success("Analysis complete. No sensitive columns suggested.")
			return
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
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
