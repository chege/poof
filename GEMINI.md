# Gemini Project Context: dbmask

This file provides comprehensive context for Gemini to understand the `dbmask` project, its architecture, and development conventions.

## Project Overview

`dbmask` is a Go-based CLI tool designed for deterministic, parallel-safe, and declarative data masking in PostgreSQL databases. It allows developers to replace sensitive data with realistic fake data using a configuration-driven approach.

### Main Technologies
- **Language**: Go (1.25.7+)
- **CLI Framework**: Cobra
- **Configuration**: HCL (HashiCorp Configuration Language)
- **Database Driver**: pgx (v5)
- **Testing**: Testcontainers-go (PostgreSQL module), Testify

### Architecture
The project follows a standard Go project layout:
- `cmd/dbmask/`: Entry points and CLI command implementations (`root.go`, `apply.go`).
- `internal/`: Private library code.
    - `config/`: HCL parsing and configuration models.
    - `db/`: PostgreSQL client wrapper and database interactions.
    - `generator/`: Data generation logic, including the `Generator` interface, a factory-based registry, and various generator types (`faker`, `constant`, `null`, `template`).
    - `masker/`: The core orchestration engine that manages parallel row processing and database updates.

### Key Concepts
- **Determinism**: Every row's masked value is seeded using `MD5(table_name + ":" + primary_key_value)`. This ensures that the same input always produces the same masked output, regardless of worker count or execution order.
- **Parallel-Safe**: A worker pool architecture separates data generation (parallelizable) from database updates (sequential within a transaction to maintain integrity).
- **Safe-by-Default**: The tool requires an explicit `allowlist` in the config or a `--force` flag to prevent accidental runs on sensitive databases.

## Building and Running

### Commands
- **Build**: `go build ./cmd/dbmask`
- **Run**: `./dbmask apply --db "postgres://user:pass@localhost:5432/dbname" --config dbmask.hcl`
- **Test**: `go test -v ./...` (Requires Docker for Testcontainers)

### Dependencies
Dependencies are managed via `go.mod`. Key dependencies include `github.com/spf13/cobra`, `github.com/hashicorp/hcl/v2`, and `github.com/jackc/pgx/v5`.

## Development Conventions

### Coding Style
- **Idiomatic Go**: Follow standard Go formatting (`gofmt`) and naming conventions.
- **Internal Package Usage**: Keep core logic in `internal/` to prevent external imports.
- **Registry Pattern**: New generators and faker providers must be explicitly registered in `internal/generator/all.go`.

### Testing Practices
- **E2E Testing**: Use `testcontainers-go` for end-to-end tests involving a real PostgreSQL instance.
- **Determinism Checks**: Always verify that masking output is consistent across different worker counts.
- **Test Fakers**: Use dedicated "test fakers" (registered in `internal/generator/test_fakers.go`) for predictable assertions in automated tests.

### Git & Workflow
- **Conventional Commits**: Use Angular prefixes for all commit messages (e.g., `feat:`, `fix:`, `docs:`, `chore:`, `test:`, `refactor:`, `perf:`).
- **Post-Implementation**: Upon completion of the Speckit workflow (Specify -> Plan -> Tasks -> Analyze -> Implement) and successful verification, the agent should:
    1.  Stage all relevant changes (excluding binaries and secrets).
    2.  Commit the changes to the feature branch using conventional commit format.
    3.  Checkout the `main` branch and fast-forward merge the feature branch.

### Contribution Guidelines
- All new features should start with a specification in the `specs/` directory following the Speckit workflow.
- Ensure all tests pass and `go mod tidy` is run before committing.
- Do not commit binaries or sensitive environment information.
