# Implementation Plan: Database Masking CLI (Poof)

**Branch**: `007-database-masking` | **Date**: Saturday, February 7, 2026 | **Spec**: [specs/007-database-masking/spec.md](specs/007-database-masking/spec.md)
**Input**: Feature specification from `specs/007-database-masking/spec.md`

## Summary

Implement the core execution modes (`plan`, `apply`, `analyze`) and safety mechanisms (retry logic, determinism) for the Poof CLI. This includes building a new analysis engine for config suggestions, enhancing the existing masking engine with retry capabilities for unique constraints, and ensuring the `plan` and `dry-run` modes provide accurate, side-effect-free previews.

## Technical Context

**Language/Version**: Go 1.25.7
**Primary Dependencies**:
- `github.com/spf13/cobra` (CLI)
- `github.com/jackc/pgx/v5` (Database Driver)
- `github.com/BurntSushi/toml` (Config Parsing)
- `golang.org/x/sync/errgroup` (Concurrency)
**Storage**: PostgreSQL 14+ (Target Database)
**Testing**: `testcontainers-go` (Integration), `stretchr/testify` (Unit)
**Target Platform**: Linux, macOS, Windows (CLI)
**Performance Goals**: >10k rows/sec masking throughput on local hardware.
**Constraints**: Zero schema modification. No superuser privileges.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Safe-by-default**: Pass. Requires explicit `apply` command and allowlist checking.
- **Deterministic**: Pass. Seeding strategy defined in spec.
- **No Magic**: Pass. Configuration is explicit TOML; `analyze` is advisory only.
- **Zero Schema Mod**: Pass. Plan explicitly forbids `ALTER`/`DROP`.

## Project Structure

### Documentation (this feature)

```text
specs/007-database-masking/
├── plan.md              # This file
├── spec.md              # Feature specification
└── checklists/          # Quality checklists
    └── requirements.md
```

### Source Code (repository root)

```text
cmd/poof/
├── main.go
├── root.go
├── apply.go        # Enhance: Add --dry-run flag handling
├── plan.go         # Enhance: Show full plan without execution
└── analyze.go      # New: Implement analyze command

internal/
├── analyze/        # New: Schema analysis and heuristic engine
│   ├── analyzer.go
│   └── rules.go
├── config/
│   ├── models.go
│   └── toml.go
├── db/
│   ├── interface.go
│   └── postgres/   # Enhance: Add schema inspection queries
├── engine/
│   ├── engine.go   # Enhance: Add retry logic, full dry-run mode
│   └── worker.go   # Enhance: Handle unique constraint retries
├── generator/
│   ├── registry.go
│   └── ...         # Enhance: Ensure all spec generators exist
└── ui/
    └── output.go   # Enhance: Reporting for retries/failures
```

**Structure Decision**: Continue with the existing modular structure. Add a new `internal/analyze` package for the advisory analysis logic to keep it separate from the core masking engine.

## Research (Phase 0)

### 1. Unique Constraint Retry Strategy
**Problem**: Randomly (or deterministically) generated values might collide with existing values in columns with `UNIQUE` constraints.
**Solution**:
- **Detection**: Catch `23505` (unique_violation) errors from PostgreSQL during `UPDATE`.
- **Retry**:
    - The `worker` will catch this error.
    - It will request a new value from the generator, passing a `retry_count` seed modifier (e.g., `seed + retry_count`).
    - Attempt the update again.
    - Repeat up to `MaxRetries` (default 10).
- **Failure**: If retries exhausted, log error and skip row (or fail job based on config).
- **Batching Impact**: Using `pgx.Batch` makes individual row retry harder. We may need to fall back to single-row updates for failed batches or use `ON CONFLICT` if we can construct it dynamically (complex without schema mod).
    - *Decision*: For v1, use single-row updates or small batches. If a batch fails, retry rows individually.

### 2. Analyze Mode Heuristics
**Goal**: Suggest columns to mask.
**Approach**:
- Query `information_schema.columns`.
- Match column names against a set of regex rules (e.g., `email`, `ssn`, `password`, `.*_name`).
- Return a list of `Table.Column` candidates with recommended generator types.
- Output: Print suggested TOML snippets to stdout.

## Data Model (Phase 1)

### Configuration Extensions (`internal/config/models.go`)
No major breaking changes.
- **Generator Params**: Ensure `Gen` struct supports all required generator types (e.g., `hash`, `counter`).
- **Retry Config**: Add `MaxRetries` to `Config` or `Table` struct (optional, can default to global constant first).

## Quickstart (Phase 1)

**Analyze**:
```bash
./poof analyze --db "postgres://..."
```

**Plan**:
```bash
./poof plan --config poof.toml
```

**Apply (Dry Run)**:
```bash
./poof apply --config poof.toml --dry-run
```

## Detailed Tasks (Phase 2 Preview)

1.  **Core**: Implement `internal/analyze` package.
    -   Schema introspection.
    -   Heuristic matching.
2.  **Engine**: Implement Retry Logic.
    -   Update `Worker` to handle `pgconn.PgError` code `23505`.
    -   Implement backoff/retry loop.
3.  **CLI**: Implement `cmd/poof/analyze.go`.
4.  **CLI**: Wire up `plan` and `apply --dry-run` to use the enhanced engine.
5.  **Generators**: Audit and implement missing generators (`hash`, `counter`).