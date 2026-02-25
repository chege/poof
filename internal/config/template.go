package config

// DefaultTemplate returns a standard TOML configuration string with optional detailed explanations.
func DefaultTemplate(explain bool) string {
	if explain {
		return `# poof configuration file (TOML)
# For more info see: https://github.com/chege/poof

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
# pk: The primary key column name. Defaults to "id".
pk = "id"

  # [[tables.columns]] section: Defines how to mask a specific column.
  [[tables.columns]]
  name = "full_name"
  # generator: Shorthand for faker providers.
  generator = "full_name"

  [[tables.columns]]
  name = "email"
  generator = "email"
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

  [[tables.columns]]
  name = "full_name"
  generator = "full_name"

  [[tables.columns]]
  name = "email"
  generator = "email"
`
}
