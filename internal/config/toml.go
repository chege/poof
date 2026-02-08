// Package config handles the loading, validation, and representation of the poof TOML configuration.
package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

var (
	// ErrMissingDSN is returned when no database connection string is provided.
	ErrMissingDSN = fmt.Errorf("database connection string is required (either in config or via --db)")
)

// LoadConfig reads and validates a TOML configuration file from the given path.
func LoadConfig(path string) (*Config, error) {
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// Expand environment variables (e.g. ${DB_PASS})
	expandedContent := os.ExpandEnv(string(content))

	var cfg Config
	// Strict validation: check for unknown keys
	metadata, err := toml.NewDecoder(strings.NewReader(expandedContent)).Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TOML: %w", err)
	}

	if len(metadata.Undecoded()) > 0 {
		return nil, fmt.Errorf("config contains unknown fields: %v", metadata.Undecoded())
	}

	// Backward compatibility: If 'database' exists but 'databases' doesn't, migrate it
	if cfg.Database != nil && len(cfg.Databases) == 0 {
		cfg.Databases = map[string]Database{
			"default": *cfg.Database,
		}
	}

	if len(cfg.Databases) == 0 {
		return nil, fmt.Errorf("config must contain at least one database definition")
	}

	// Run declarative validation
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 500
	}

	// Apply Sane Defaults
	for i := range cfg.Tables {
		if cfg.Tables[i].PK == "" {
			cfg.Tables[i].PK = "id"
		}
		for j := range cfg.Tables[i].Columns {
			col := &cfg.Tables[i].Columns[j]
			// If shorthand 'generator' is used, populate the 'Gen' struct
			if col.Generator != "" && col.Gen.Type == "" {
				col.Gen = Gen{
					Type:     "faker",
					Provider: col.Generator,
				}
			}
		}
	}

	for env, dbEnv := range cfg.Databases {
		if dbEnv.DSN == "" {
			// Try to build DSN from standard PG env vars
			host := os.Getenv("PGHOST")
			if host == "" {
				host = "localhost"
			}
			user := os.Getenv("PGUSER")
			pass := os.Getenv("PGPASSWORD")
			dbname := os.Getenv("PGDATABASE")
			port := os.Getenv("PGPORT")
			if port == "" {
				port = "5432"
			}

			if user != "" && dbname != "" {
				dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbname)
				dbEnv.DSN = dsn
				cfg.Databases[env] = dbEnv
			}
		}
	}

	return &cfg, nil
}

// GetDatabase returns the connection details for a specific environment.
func (c *Config) GetDatabase(env string) (Database, error) {
	if env != "" {
		db, ok := c.Databases[env]
		if !ok {
			return Database{}, fmt.Errorf("database environment %q not found", env)
		}
		return db, nil
	}

	// Try to find the default
	for _, db := range c.Databases {
		if db.Default {
			return db, nil
		}
	}

	// Fallback to "default" named block
	if db, ok := c.Databases["default"]; ok {
		return db, nil
	}

	// Last resort: return the first one
	for _, db := range c.Databases {
		return db, nil
	}

	return Database{}, fmt.Errorf("no database environments configured")
}

// Save writes the configuration to the given writer in TOML format.
func (c *Config) Save(w io.Writer) error {
	return toml.NewEncoder(w).Encode(c)
}
