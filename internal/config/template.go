package config

// DefaultTemplate returns a standard TOML configuration string with optional detailed explanations.
func DefaultTemplate(explain bool) string {
	if explain {
		return `# poof configuration file (TOML)
# For more info see: https://github.com/christopher/poof

# [databases] section: Connection details for multiple environments.
[databases.local]
# dsn: The Data Source Name for connecting to the database.
# Supporting environment variables is highly recommended for security.
dsn = "postgres://localhost:5432/my_app?sslmode=disable"
default = true

[databases.staging]
dsn = "postgres://staging_user:${STAGING_DB_PASS}@staging-host:5432/my_app"

# [safety] section: Guardrails for irreversible masking.
[safety]
# allowed_db_names: List of databases where masking is allowed without --force.
allowed_db_names = ["my_app", "my_app_staging"]
# salt: A secret string to ensure different environments mask data differently.
salt = "${POOF_SALT}"

# locale: Default language for faker generators (e.g. en_US, de_DE).
locale = "en_US"

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
  # type: Generator type (faker, template, constant, null, hash, counter).
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

	return `[databases.local]
dsn = "postgres://localhost:5432/testdb?sslmode=disable"
default = true

[safety]
allowed_db_names = ["testdb"]
salt = "secret-salt-change-me"

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
}