package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/christopher/poof/internal/db"
	"github.com/christopher/poof/internal/ui"
)

// GeneratorValidator defines the interface for validating generator parameters.
// This interface allows the config package to validate generators without
// importing the generator package (avoiding import cycles).
type GeneratorValidator interface {
	ValidateGen(locale string, gen Gen, sqlDataType string) error
}

// ValidateStatic performs semantic validation that does not require a database connection.
func (c *Config) ValidateStatic(gv GeneratorValidator) error {
	// 1. Basic database validation
	for env, dbEnv := range c.Databases {
		if dbEnv.DSN == "" {
			return fmt.Errorf("environment %q: dsn is required", env)
		}
	}

	for _, table := range c.Tables {

		pk := table.PK

		if pk == "" {

			pk = "id" // Default

		}

		for _, col := range table.Columns {

			if col.Name == pk {

				return fmt.Errorf("table %q: cannot mask the primary key column %q", table.Name, pk)

			}

			gen := col.Gen

			if col.Generator != "" && gen.Type == "" {

				gen = Gen{Type: "faker", Provider: col.Generator}

			}

			if err := gv.ValidateGen(c.Locale, gen, ""); err != nil {

				return fmt.Errorf("table %q, column %q: %w", table.Name, col.Name, err)

			}

		}

	}

	return nil

}

// ValidateDatabase performs semantic validation against a live database schema.
func (c *Config) ValidateDatabase(ctx context.Context, client db.DB, gv GeneratorValidator) error {
	dbName, err := client.GetDatabaseName(ctx)
	if err != nil {
		return fmt.Errorf("failed to get database name: %w", err)
	}

	// 1. Production Readiness Check

	isProd := strings.Contains(strings.ToLower(dbName), "prod") ||

		strings.Contains(strings.ToLower(dbName), "live")

	if isProd && c.Safety.Salt == "" {

		return fmt.Errorf("production-like database %q detected but no global salt is configured in [safety]", dbName)

	}

	// 2. Schema Verification & Type Check

	for _, tableCfg := range c.Tables {

		cols, err := client.GetTableColumns(ctx, tableCfg.Name)

		if err != nil {

			return fmt.Errorf("failed to inspect table %q: %w", tableCfg.Name, err)

		}

		if len(cols) == 0 {

			return fmt.Errorf("table %q not found in database %q", tableCfg.Name, dbName)

		}

		colMap := make(map[string]db.ColumnInfo)

		for _, col := range cols {

			colMap[col.Name] = col

		}

		for _, colCfg := range tableCfg.Columns {

			info, exists := colMap[colCfg.Name]

			if !exists {

				return fmt.Errorf("column %q not found in table %q", colCfg.Name, tableCfg.Name)

			}

			if info.IsForeignKey {

				ui.Warning("Table %q: Column %q is a FOREIGN KEY. Masking it may break relational integrity.", tableCfg.Name, colCfg.Name)

			}

			// Perform type compatibility check

			gen := colCfg.Gen

			if colCfg.Generator != "" && gen.Type == "" {

				gen = Gen{Type: "faker", Provider: colCfg.Generator}

			}

			if err := gv.ValidateGen(c.Locale, gen, info.DataType); err != nil {

				return fmt.Errorf("table %q, column %q: %w", tableCfg.Name, colCfg.Name, err)

			}

		}

	}

	return nil

}
