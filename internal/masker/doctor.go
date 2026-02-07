package masker

import (
	"context"

	"github.com/christopher/masker/internal/config"
	"github.com/christopher/masker/internal/db"
	_ "github.com/christopher/masker/internal/db/postgres"
	"github.com/christopher/masker/internal/generator"
	"github.com/christopher/masker/internal/ui"
)

func CheckReadiness(ctx context.Context, configPath string, dbConnStr string) bool {
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
	client, err := db.Connect(ctx, dbConnStr)
	if err != nil {
		ui.Error("Database connectivity: %v", err)
		success = false
	} else {
		defer client.Close()
		dbName, err := client.GetDatabaseName(ctx)
		if err != nil {
			ui.Error("Database name retrieval: %v", err)
			success = false
		} else {
			ui.Success("Connected to database: %s", dbName)

			// 3. Safety Check
			if cfg != nil {
				if cfg.IsAllowed(dbName) {
					ui.Success("Database is in the allowlist")
				} else {
					ui.Warning("Database %q is not in the allowlist (requires --force)", dbName)
				}
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
