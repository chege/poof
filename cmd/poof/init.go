// Package main provides the CLI entry point for poof.
package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/christopher/poof/internal/config"
	"github.com/christopher/poof/internal/ui"
)

var explain bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new poof.toml configuration file",
	Run: func(_ *cobra.Command, _ []string) {
		path := "poof.toml"
		if _, err := os.Stat(path); err == nil {
			ui.Error("File %s already exists. Refusing to overwrite.", path)
			os.Exit(1)
		}

		content := config.DefaultTemplate(explain)
		// #nosec G306 -- configuration file is readable by the owner.
		err := os.WriteFile(path, []byte(content), 0600)
		if err != nil {
			ui.Error("Failed to write %s: %v", path, err)
			os.Exit(1)
		}

		ui.Success("Initialized %s", path)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&explain, "explain", false, "Include detailed inline explanations")
}
