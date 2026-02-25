// Package main provides the CLI entry point for poof.
package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/chege/poof/internal/config"
	"github.com/chege/poof/internal/ui"
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
		ui.Info("Next steps:")
		ui.Info("  1. Edit %s to set your database connection.", path)
		ui.Info("  2. Run 'poof analyze' to automatically suggest masking rules.")
		ui.Info("  3. Run 'poof plan' to preview the transformation.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&explain, "explain", false, "Include detailed inline explanations")
}
