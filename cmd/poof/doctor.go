// Package main provides the CLI entry point for poof.
package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/christopher/poof/internal/generator"
	"github.com/christopher/poof/internal/ui"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run pre-flight readiness checks",
	Long: `Performs a series of diagnostic checks to ensure the environment
is correctly configured for masking.
Checks include database connectivity, configuration semantic validity,
and schema existence.`,
	Run: func(cmd *cobra.Command, _ []string) {
		ui.HandleExit(runDoctor(cmd.Context()))
	},
}

func runDoctor(ctx context.Context) error {
	cli, err := LoadResources(ctx)
	if err != nil {
		return err
	}
	defer cli.DB.Close()

	ui.Info("Running pre-flight checks...")

	// 1. Static Validation
	generator.RegisterAll()
	if err := cli.Config.ValidateStatic(&generator.Validator{}); err != nil {
		ui.Error("Static validation: %v", err)
		return ui.WrapError(ui.ErrConfig, err)
	}
	ui.Success("Configuration semantics are valid")

	// 2. Database Validation
	if err := cli.Config.ValidateDatabase(ctx, cli.DB); err != nil {
		ui.Error("Database validation: %v", err)
		return ui.WrapError(ui.ErrConfig, err)
	}
	ui.Success("Database schema and safety checks passed")

	ui.Success("Everything looks good! You are ready to mask.")
	return nil
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
