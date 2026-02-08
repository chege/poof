// Package main provides the CLI entry point for poof.
package main

import (
	"context"

	"github.com/spf13/cobra"

	_ "github.com/christopher/poof/internal/db/postgres"
	"github.com/christopher/poof/internal/engine"
	"github.com/christopher/poof/internal/ui"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check if the environment is ready for masking",
	Run: func(_ *cobra.Command, _ []string) {
		ui.HandleExit(runDoctor())
	},
}

func runDoctor() error {
	ctx := context.Background()

	cli, err := LoadResources(ctx)
	if err != nil {
		return err
	}
	defer cli.DB.Close()

	if !engine.CheckReadiness(ctx, configPath, cli.DB) {
		return ui.WrapError(ui.ErrConnection, nil)
	}
	ui.Success("Everything looks good! You are ready to mask.")
	return nil
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
