// Package main provides the CLI entry point for poof.
package main

import (
	"github.com/spf13/cobra"

	"github.com/christopher/poof/internal/config"
	"github.com/christopher/poof/internal/ui"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the TOML configuration file",
	Run: func(_ *cobra.Command, _ []string) {
		ui.HandleExit(runValidate())
	},
}

func runValidate() error {
	_, err := config.LoadConfig(configPath)
	if err != nil {
		return ui.WrapError(ui.ErrConfig, err)
	}
	ui.Success("Configuration is valid.")
	return nil
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
