package config

type Config struct {
	Database Database `toml:"database" validate:"required"`
	Safety   Safety   `toml:"safety" validate:"required"`
	Tables   []Table  `toml:"tables" validate:"required,dive,min=1"`
}

type Database struct {
	DSN string `toml:"dsn" validate:"required"`
}

type Safety struct {
	AllowedDBNames []string `toml:"allowed_db_names" validate:"required,min=1"`
}

type Table struct {
	Name    string   `toml:"name" validate:"required"`
	PK      string   `toml:"pk" validate:"required"`
	Source  *Source  `toml:"source" validate:"omitempty,dive"`
	Columns []Column `toml:"columns" validate:"required,dive,min=1"`
}

type Column struct {
	Name string `toml:"name" validate:"required"`
	Gen  Gen    `toml:"gen" validate:"required"`
}

type Gen struct {
	Type     string            `toml:"type" validate:"required,oneof=faker template constant null"`
	Provider string            `toml:"provider" validate:"required_if=Type faker"`
	Value    string            `toml:"value" validate:"required_if=Type constant"`
	Template string            `toml:"template" validate:"required_if=Type template"`
	Params   map[string]string `toml:"params" validate:"omitempty"`
}

type Source struct {
	Type   string            `toml:"type" validate:"required,oneof=table view query"`
	Name   string            `toml:"name" validate:"required_if=Type view"`
	SQL    string            `toml:"sql" validate:"required_if=Type query"`
	Params map[string]string `toml:"params" validate:"omitempty"`
}

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
