// Package main provides the CLI entry point for poof.
package main

import (
	"os"

	"github.com/christopher/poof/internal/config"
	"github.com/christopher/poof/internal/ui"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the TOML configuration file",
	Run: func(_ *cobra.Command, _ []string) {
		_, err := config.LoadConfig(configPath)
		if err != nil {
			ui.Error("Validation failed: %v", err)
			os.Exit(1)
		}
		ui.Success("Configuration is valid.")
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
