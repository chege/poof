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

var (
	force bool
	yes   bool
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply data masking rules to the database",
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

		dbName, err := client.GetDatabaseName(ctx)
		if err != nil {
			ui.Error("DB name error: %v", err)
			os.Exit(1)
		}

		if !cfg.IsAllowed(dbName) && !force {
			ui.Error("Database %q is not in the allowlist and --force was not provided.", dbName)
			os.Exit(1)
		}

		if !yes {
			ui.Info("Running plan before apply...")
			engine := masker.NewEngine(client, cfg, workers)
			engine.DryRun = true
			_, err := engine.Apply(ctx)
			if err != nil {
				ui.Error("Pre-apply plan failed: %v", err)
				os.Exit(1)
			}
			ui.Success("Plan summary generated.")
		}

		ui.Info("Applying masking...")
		engine := masker.NewEngine(client, cfg, workers)
		_, err = engine.Apply(ctx)
		if err != nil {
			ui.Error("Apply failed: %v", err)
			os.Exit(1)
		}

		ui.Success("Masking completed successfully.")
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().BoolVar(&force, "force", false, "Force apply even if database is not in allowlist")
	applyCmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip plan summary and proceed immediately")
}
