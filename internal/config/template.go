package config

func DefaultTemplate(explain bool) string {
	content := `[database]
dsn = "postgres://user:pass@localhost:5432/testdb"

[safety]
allowed_db_names = ["testdb"]

[[tables]]
name = "users"
pk = "id"

  [[tables.columns]]
  name = "full_name"
  [tables.columns.gen]
  type = "faker"
  provider = "full_name"

  [[tables.columns]]
  name = "email"
  [tables.columns.gen]
  type = "faker"
  provider = "email"
`
	if explain {
		content = `# dbmask configuration file (TOML)
# For more info see: https://github.com/christopher/masker

# [database] section: Connection details.
[database]
# dsn: The Data Source Name for connecting to the database.
dsn = "postgres://user:pass@localhost:5432/testdb"

# [safety] section: Guardrails for irreversible masking.
[safety]
# allowed_db_names: List of databases where masking is allowed without --force.
allowed_db_names = ["testdb"]

# [[tables]] section: Define rules for a table (can be multiple).
[[tables]]
name = "users"
# pk: The primary key column name. REQUIRED for deterministic masking.
pk = "id"

  # [[tables.columns]] section: Defines how to mask a specific column.
  [[tables.columns]]
  name = "full_name"
  # [tables.columns.gen] section: Generator configuration.
  [tables.columns.gen]
  # type: Generator type (faker, template, constant, null).
  type = "faker"
  # provider: Specific data type for faker.
  provider = "full_name"

  [[tables.columns]]
  name = "email"
  [tables.columns.gen]
  type = "faker"
  provider = "email"
`
	}
	return content
}
