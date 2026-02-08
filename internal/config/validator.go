package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/christopher/poof/internal/db"
)

// GeneratorValidator defines the interface for validating generator parameters.
// This interface allows the config package to validate generators without
// importing the generator package (avoiding import cycles).
type GeneratorValidator interface {
	ValidateGen(locale string, gen Gen) error
}

// ValidateStatic performs semantic validation that does not require a database connection.
func (c *Config) ValidateStatic(gv GeneratorValidator) error {
	for _, table := range c.Tables {
		for _, col := range table.Columns {
			if err := gv.ValidateGen(c.Locale, col.Gen); err != nil {
				return fmt.Errorf("table %q, column %q: %w", table.Name, col.Name, err)
			}
		}
	}
	return nil
}

// ValidateDatabase performs semantic validation against a live database schema.
func (c *Config) ValidateDatabase(ctx context.Context, client db.DB) error {
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

	// 2. Schema Verification
	for _, tableCfg := range c.Tables {
		cols, err := client.GetTableColumns(ctx, tableCfg.Name)
		if err != nil {
			return fmt.Errorf("failed to inspect table %q: %w", tableCfg.Name, err)
		}
		if len(cols) == 0 {
			return fmt.Errorf("table %q not found in database %q", tableCfg.Name, dbName)
		}

		colMap := make(map[string]bool)
		for _, col := range cols {
			colMap[col.Name] = true
		}

		for _, colCfg := range tableCfg.Columns {
			if !colMap[colCfg.Name] {
				return fmt.Errorf("column %q not found in table %q", colCfg.Name, tableCfg.Name)
			}
		}
	}

	return nil
}
