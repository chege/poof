package main

import (
	"context"
	"os"

	"github.com/christopher/masker/internal/masker"
	"github.com/christopher/masker/internal/ui"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check if the environment is ready for masking",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		if dbConnStr == "" {
			ui.Error("Database connection string is required (--db)")
			os.Exit(1)
		}
		if !masker.CheckReadiness(ctx, configPath, dbConnStr) {
			ui.Error("Doctor found issues. Please fix them before running apply.")
			os.Exit(1)
		}
		ui.Success("Everything looks good! You are ready to mask.")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
