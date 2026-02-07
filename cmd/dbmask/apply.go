package main

import (
	"context"
	"os"

	"github.com/christopher/masker/internal/config"
	"github.com/christopher/masker/internal/db"
	_ "github.com/christopher/masker/internal/db/postgres"
	"github.com/christopher/masker/internal/generator"
	"github.com/christopher/masker/internal/masker"
	"github.com/christopher/masker/internal/ui"
	"github.com/spf13/cobra"
)

var (
	force  bool
	yes    bool
	dryRun bool
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
			ui.Error("DB connection error: %v", err)
			os.Exit(1)
		}
		defer client.Close()

		dbName, err := client.GetDatabaseName(ctx)
		if err != nil {
			ui.Error("DB name error: %v", err)
			os.Exit(1)
		}

		if !cfg.IsAllowed(dbName) && !force {
			ui.Error("Database %q is not in the allowed_db_names list and --force was not provided.", dbName)
			os.Exit(1)
		}

		if !yes && !dryRun {
			ui.Info("Running plan before apply...")
			engine := masker.NewEngine(client, cfg, workers)
			engine.DryRun = true
			_, err := engine.Apply(ctx)
			if err != nil {
				ui.Error("Pre-apply plan failed: %v", err)
				os.Exit(1)
			}
			ui.Success("Pre-flight plan completed.")
		}

		if dryRun {
			ui.Info("Applying masking (DRY RUN)...")
		} else {
			ui.Info("Applying masking...")
		}

		engine := masker.NewEngine(client, cfg, workers)
		engine.DryRun = dryRun
		_, err = engine.Apply(ctx)
		if err != nil {
			ui.Error("Apply failed: %v", err)
			os.Exit(1)
		}

		if dryRun {
			ui.Success("Dry run completed successfully. No changes were made.")
		} else {
			ui.Success("Masking completed successfully.")
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().BoolVar(&force, "force", false, "Force apply even if database is not in allowlist")
	applyCmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip plan summary and proceed immediately")
	applyCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Execute masking but do not commit changes")
}
