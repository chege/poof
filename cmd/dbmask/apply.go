package main

import (
	"context"
	"fmt"
	"log"

	"github.com/christopher/masker/internal/config"
	"github.com/christopher/masker/internal/db"
	"github.com/christopher/masker/internal/generator"
	"github.com/christopher/masker/internal/masker"
	"github.com/spf13/cobra"
)

var (
	configPath string
	dbConnStr  string
	force      bool
	workers    int
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply data masking rules to the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		generator.RegisterAll()

		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			return err
		}

		ctx := context.Background()
		client, err := db.NewClient(ctx, dbConnStr)
		if err != nil {
			return err
		}
		defer client.Close()

		dbName, err := client.GetDatabaseName(ctx)
		if err != nil {
			return err
		}

		if !cfg.IsAllowed(dbName) && !force {
			return fmt.Errorf("database %q is not in the allowlist and --force was not provided", dbName)
		}

		engine := masker.NewEngine(client, cfg, workers)
		if err := engine.Apply(ctx); err != nil {
			return err
		}

		log.Println("Masking completed successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringVarP(&configPath, "config", "c", "dbmask.hcl", "Path to HCL config file")
	applyCmd.Flags().StringVarP(&dbConnStr, "db", "d", "", "PostgreSQL connection string (e.g. postgres://user:pass@localhost:5432/dbname)")
	applyCmd.Flags().BoolVar(&force, "force", false, "Force apply even if database is not in allowlist")
	applyCmd.Flags().IntVarP(&workers, "workers", "w", 4, "Number of parallel workers")

	applyCmd.MarkFlagRequired("db")
}
