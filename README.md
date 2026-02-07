# poof

`poof` is a PostgreSQL data masking tool designed to be deterministic, parallel-safe, and declarative.

## Features

- **Declarative HCL Configuration**: Define masking rules in a human-readable format.
- **Deterministic**: Masking is based on a seed derived from the table name and primary key, ensuring consistent results across runs.
- **Parallel-Safe**: Uses a worker pool for efficient data generation without affecting determinism.
- **Safe-by-Default**: Refuses to run on databases not in the allowlist unless `--force` is provided.
- **Taskfile Integration**: Standardized development and orchestration via Taskfile.dev.
- **Extensible**: Easily add new generators or faker providers in Go.

## Development Workflow

This project uses `task` (Taskfile.dev) for common operations:

- **Build**: `task build`
- **Test**: `task test`
- **Lint/Check**: `task check` (runs fmt, vet, and test)
- **Apply Masking**: `task apply -- DB_URL="your-db-url" CONFIG_PATH="your-config.hcl"`

## Usage

Create a `poof.hcl` file:

```hcl
allowlist = ["testdb"]

table "users" {
  pk = "id"

  column "full_name" {
    gen "faker" {
      provider = "full_name"
    }
  }

  column "username" {
    gen "faker" {
      provider = "username"
    }
  }
}
```

Run the masking tool:

```bash
./poof apply --db "postgres://user:pass@localhost:5432/testdb" --config poof.hcl
```

## Supported Generators

- `faker`: Uses a provider to generate fake data. Supported providers:
    - `first_name`
    - `last_name`
    - `full_name`
    - `username`
    - `email`
    - `company_name`
    - `phone_number`
    - `ipv4_address`
    - `short_text`
- `constant`: Returns a fixed value.
- `null`: Sets the column to NULL.
- `template`: Uses Go templates to generate values.

## Data Producers

You can control how rows are selected for masking using the `source` block:

### Table (Default)
```hcl
table "users" {
  pk = "id"
  # Omitted source defaults to table scan
}
```

### View
```hcl
table "users" {
  pk = "id"
  source "view" {
    name = "active_users_view"
  }
}
```

### Custom Query
```hcl
table "users" {
  pk = "id"
  source "query" {
    sql = "SELECT id FROM users WHERE active = true ORDER BY id"
  }
}
```
*Note: Custom queries MUST include `ORDER BY pk` to ensure determinism.*

## Testing

```bash
task test
```