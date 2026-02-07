# dbmask

`dbmask` is a PostgreSQL data masking tool designed to be deterministic, parallel-safe, and declarative.

## Features

- **Declarative HCL Configuration**: Define masking rules in a human-readable format.
- **Deterministic**: Masking is based on a seed derived from the table name and primary key, ensuring consistent results across runs.
- **Parallel-Safe**: Uses a worker pool for efficient data generation without affecting determinism.
- **Safe-by-Default**: Refuses to run on databases not in the allowlist unless `--force` is provided.
- **Extensible**: Easily add new generators or faker providers in Go.

## Installation

```bash
go build ./cmd/dbmask
```

## Usage

Create a `dbmask.hcl` file:

```hcl
allowlist = ["testdb"]

table "users" {
  pk = "id"

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
```

Run the masking tool:

```bash
./dbmask apply --db "postgres://user:pass@localhost:5432/testdb" --config dbmask.hcl
```

## Generators

- `faker`: Uses a provider to generate fake data (e.g., `first_name`, `email`, `full_name`, `company`).
- `constant`: Returns a fixed value.
- `null`: Sets the column to NULL.
- `template`: Uses Go templates to generate values.

## Testing

```bash
go test ./...
```
