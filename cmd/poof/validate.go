// Package main provides the CLI entry point for poof.
package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/christopher/poof/internal/generator"
	"github.com/christopher/poof/internal/ui"
)

var (
	dbCheck bool
	strict  bool
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the TOML configuration file",
	Run: func(cmd *cobra.Command, _ []string) {
		ui.HandleExit(runValidate(cmd.Context()))
	},
}

func runValidate(ctx context.Context) error {
	cli, err := LoadResources(ctx)
	if err != nil {
		return err
	}
	defer cli.DB.Close()

	// Level 1 & 2: Syntax + Static Semantic
	ui.Info("Running static validation...")
	generator.RegisterAll()
	if err := cli.Config.ValidateStatic(&generator.Validator{}); err != nil {
		return ui.WrapError(ui.ErrConfig, err)
	}

	// Level 3: Database Semantic
	if dbCheck {
		ui.Info("Running database schema validation...")
		if err := cli.Config.ValidateDatabase(ctx, cli.DB, &generator.Validator{}); err != nil {
			return ui.WrapError(ui.ErrConfig, err)
		}
	}

	ui.Success("Configuration is valid.")
	return nil
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().BoolVar(&dbCheck, "db-check", false, "Verify configuration against the live database schema")
	validateCmd.Flags().BoolVar(&strict, "strict", false, "Treat all warnings as errors")
}
