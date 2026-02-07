package config

type Config struct {
	Database Database `toml:"database"`
	Safety   Safety   `toml:"safety"`
	Tables   []Table  `toml:"tables"`
}

type Database struct {
	DSN string `toml:"dsn"`
}

type Safety struct {
	AllowedDBNames []string `toml:"allowed_db_names"`
}

type Table struct {
	Name    string   `toml:"name"`
	PK      string   `toml:"pk"`
	Source  *Source  `toml:"source"`
	Columns []Column `toml:"columns"`
}

type Column struct {
	Name string `toml:"name"`
	Gen  Gen    `toml:"gen"`
}

type Gen struct {
	Type     string            `toml:"type"`
	Provider string            `toml:"provider"`
	Value    string            `toml:"value"`
	Template string            `toml:"template"`
	Params   map[string]string `toml:"params"`
}

type Source struct {
	Type   string            `toml:"type"`
	Name   string            `toml:"name"`
	SQL    string            `toml:"sql"`
	Params map[string]string `toml:"params"`
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
