// Package config defines the data models for poof configuration.
package config

// Config represents the top-level poof configuration structure.
type Config struct {
	Database  Database            `toml:"database"` // Deprecated: use Databases instead
	Databases map[string]Database `toml:"databases"`
	Locale    string              `toml:"locale"`
	Safety    Safety              `toml:"safety" validate:"required"`
	Tables    []Table             `toml:"tables" validate:"required,dive,min=1"`
}

// Database contains the connection parameters for the target database.
type Database struct {
	DSN     string `toml:"dsn" validate:"required"`
	Default bool   `toml:"default"`
}

// Safety defines the guardrails for irreversible masking operations.
type Safety struct {
	AllowedDBNames []string `toml:"allowed_db_names" validate:"required,min=1"`
}

// Table defines masking rules and source selection for a specific database table.
type Table struct {
	Source  *Source  `toml:"source" validate:"omitempty,dive"`
	Name    string   `toml:"name" validate:"required"`
	PK      string   `toml:"pk" validate:"required"`
	Columns []Column `toml:"columns" validate:"required,dive,min=1"`
}

// Column defines which generator to use for a specific database column.
type Column struct {
	Name string `toml:"name" validate:"required"`
	Gen  Gen    `toml:"gen" validate:"required"`
}

// Gen defines the generator type and its specific configuration parameters.
type Gen struct {
	Params   map[string]string `toml:"params" validate:"omitempty"`
	Type     string            `toml:"type" validate:"required,oneof=faker template constant null"`
	Provider string            `toml:"provider" validate:"required_if=Type faker"`
	Value    string            `toml:"value" validate:"required_if=Type constant"`
	Template string            `toml:"template" validate:"required_if=Type template"`
}

// Source defines an alternative row selection strategy (e.g. view or custom query).
type Source struct {
	Params map[string]string `toml:"params" validate:"omitempty"`
	Type   string            `toml:"type" validate:"required,oneof=table view query"`
	Name   string            `toml:"name" validate:"required_if=Type view"`
	SQL    string            `toml:"sql" validate:"required_if=Type query"`
}

// IsAllowed returns true if the given database name is in the list of allowed databases.
func (c *Config) IsAllowed(dbName string) bool {
	if len(c.Safety.AllowedDBNames) == 0 {
		return false
	}
	for _, allowed := range c.Safety.AllowedDBNames {
		if allowed == dbName {
			return true
		}
	}
	return false
}
