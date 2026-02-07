package main

import (
	"os"

	"github.com/christopher/masker/internal/config"
	"github.com/christopher/masker/internal/ui"
	"github.com/spf13/cobra"
)

var explain bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new dbmask.toml configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		path := "dbmask.toml"
		if _, err := os.Stat(path); err == nil {
			ui.Error("File %s already exists. Refusing to overwrite.", path)
			os.Exit(1)
		}

		content := config.DefaultTemplate(explain)
		err := os.WriteFile(path, []byte(content), 0644)
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
