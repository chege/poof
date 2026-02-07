package config

func DefaultTemplate(explain bool) string {
	content := `allowlist = ["testdb", "production"]

table "users" {
  pk = "id" # Primary key is REQUIRED for deterministic masking

  column "full_name" {
    gen "faker" {
      provider = "full_name"
    }
  }

  column "email" {
    gen "faker" {
      provider = "email"
    }
  }
}
`
	if explain {
		content = `# dbmask configuration file
# For more info see: https://github.com/christopher/masker

# allowlist: A list of database names where masking is allowed to run without --force.
allowlist = ["testdb", "production"]

# table: Defines masking rules for a specific table.
table "users" {
  # pk: The primary key column name. Used to seed the random generator for each row.
  pk = "id"

  # column: Defines how to mask a specific column.
  column "full_name" {
    # gen: Specifies the generator type and its configuration.
    gen "faker" {
      # provider: The specific faker data type to use.
      provider = "full_name"
    }
  }

  column "email" {
    gen "faker" {
      provider = "email"
    }
  }
}
`
	}
	return content
}
