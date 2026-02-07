package main

import (
	"os"

	"github.com/christopher/masker/internal/config"
	"github.com/christopher/masker/internal/ui"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the TOML configuration file",
	Run: func(cmd *cobra.Command, args []string) {
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
