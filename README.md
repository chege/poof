# poof

`poof` is a PostgreSQL data masking tool designed to be deterministic, parallel-safe, and declarative.

## Try poof in 2 minutes

Experience `poof` immediately without any setup.

**Requirements:**
- Docker
- Go (1.25.7+)

**Run the demo:**
```bash
task demo
```

**What this does:**
1. Spins up a local PostgreSQL container.
2. Loads sample user data.
3. Runs `poof plan` to show you exactly how the data will be masked.

To clean up: `task demo:clean`

## Features

- **Declarative TOML Configuration**: Define masking rules in a human-readable format.
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

Create a `poof.toml` file:

```toml
[database]
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
  name = "username"
  [tables.columns.gen]
  type = "faker"
  provider = "username"
```

Run the masking tool:

```bash
./poof apply --config poof.toml
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
```toml
[[tables]]
name = "users"
pk = "id"
# Omitted source defaults to table scan
```

### View
```toml
[[tables]]
name = "users"
pk = "id"
[tables.source]
type = "view"
name = "active_users_view"
```

### Custom Query
```toml
[[tables]]
name = "users"
pk = "id"
[tables.source]
type = "query"
sql = "SELECT id FROM users WHERE active = true ORDER BY id"
```
*Note: Custom queries MUST include `ORDER BY pk` to ensure determinism.*

## Testing

```bash
task test
```