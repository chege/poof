// Package engine coordinates the data masking process.
package engine

import (
	"context"

	"github.com/christopher/poof/internal/config"
	"github.com/christopher/poof/internal/db"
	_ "github.com/christopher/poof/internal/db/postgres" // Registered for side-effect database discovery.
	"github.com/christopher/poof/internal/generator"
	"github.com/christopher/poof/internal/ui"
)

// CheckReadiness performs a series of environment and configuration checks.
func CheckReadiness(ctx context.Context, configPath string, client db.DB) bool {
	success := true

	// 1. Config Check
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		ui.Error("Config validation: %v", err)
		success = false
	} else {
		ui.Success("Config is valid")
	}

	// 2. Database Check
	dbName, err := client.GetDatabaseName(ctx)
	if err != nil {
		ui.Error("Database name retrieval: %v", err)
		success = false
	} else {
		ui.Success("Connected to database: %s", dbName)

		// Check Job State
		state, err := client.GetJobState(ctx)
		if err == nil && (state == "STARTED" || state == "FAILED") {
			ui.Warning("Database masking state is %q. A previous run may have failed or was interrupted.", state)
		} else if state == "COMPLETED" {
			ui.Success("Previous masking run was completed successfully")
		}

		// 3. Safety Check

		if cfg != nil {
			if cfg.IsAllowed(dbName) {
				ui.Success("Database is in the allowed_db_names list")
			} else {
				ui.Warning("Database %q is not in the allowed_db_names list (requires --force)", dbName)
			}
		}
	}

	// 4. Providers Check
	if cfg != nil {
		generator.RegisterAll()
		for _, table := range cfg.Tables {
			for _, col := range table.Columns {
				_, err := generator.NewGenerator(col.Gen)
				if err != nil {
					ui.Error("Generator for %s.%s: %v", table.Name, col.Name, err)
					success = false
				}
			}
		}
		if success {
			ui.Success("All generators and providers are available")
		}
	}

	return success
}
