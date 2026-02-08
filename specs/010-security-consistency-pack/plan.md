# Implementation Plan: Security & Consistency Pack

**Branch**: `010-security-consistency-pack` | **Date**: Saturday, February 7, 2026 | **Spec**: [specs/010-security-consistency-pack/spec.md](specs/010-security-consistency-pack/spec.md)
**Input**: Feature specification from `specs/010-security-consistency-pack/spec.md`

## Summary

This plan implements five key infrastructure improvements: cross-table consistency via value-based seeding, graceful shutdown handling, global seeding salt for environment isolation, incremental masking filters, and JSON export for schema analysis.

## Technical Context

**Language/Version**: Go 1.25.7
**Primary Dependencies**: 
- `github.com/jackc/pgx/v5`
- `github.com/spf13/cobra`
- `os/signal` (Graceful shutdown)
- `encoding/json` (Export)

## Project Structure

### Documentation

```text
specs/010-security-consistency-pack/
├── plan.md              # This file
├── spec.md              # Feature specification
└── checklists/
    └── requirements.md
```

### Source Code

```text
cmd/poof/
├── root.go             # Update: Implement SIGINT handling
└── analyze.go          # Update: Add --json flag

internal/
├── config/
│   └── models.go       # Update: Add Salt, Filter, and SeedBy fields
├── engine/
│   └── engine.go       # Update: Respect context cancellation, pass original values
├── generator/
│   ├── context.go      # Update: Add Salt and OriginalValue to context
│   └── interface.go    # Update: Add OriginalValue to RowContext
└── producer/
    └── table.go        # Update: Implement SQL filters
```

## Research & Recommendations

### 1. Seeding Algorithm Refactor
**Current**: `MD5(table + ":" + pk)`
**New**: `MD5(salt + ":" + identifier)`
- If `seed_by = "pk"` (default): `identifier = table + ":" + pk`
- If `seed_by = "value"`: `identifier = original_value`
**Recommendation**: Always include the global `salt` if provided. If `salt` is empty, seeding remains backward compatible.

### 2. Graceful Shutdown
**Implementation**: 
- In `cmd/poof/root.go`, create a cancellable context using `signal.NotifyContext`.
- Pass this context down through `LoadResources`, `Engine.Apply`, and `maskTable`.
- The `Engine` already uses `errgroup.WithContext(ctx)`, so workers will stop naturally when the context is cancelled.
- **Recommendation**: Ensure `tx.Rollback()` is called on cancellation to avoid partial table commits.

### 3. Incremental Masking (SQL Filters)
**Implementation**:
- Update `producer.FetchRows` to append a `WHERE` clause if `tableCfg.Source.Filter` is present.
- **Recommendation**: Provide a clear warning if the user provides a filter without an `ORDER BY` (though we already append `ORDER BY PK`).

## Data Model

### Configuration Extensions (`internal/config/models.go`)
```go
type Safety struct {
    // ...
    Salt string `toml:"salt"`
}

type Column struct {
    // ...
    SeedBy string `toml:"seed_by"` // "pk" or "value"
}

type Source struct {
    // ...
    Filter string `toml:"filter"` // SQL WHERE clause
}
```

## Detailed Tasks

1.  **Foundations**: Update `config.Config` models and `generator.RowContext`.
2.  **Security**: Update `generator.NewRowContext` to incorporate global salt.
3.  **Consistency**: Modify `Engine` to read original values and support `seed_by = "value"`.
4.  **Performance**: Update `producer` to respect SQL `filter` clauses.
5.  **Robustness**: Implement `signal.NotifyContext` in `root.go`.
6.  **Automation**: Add `--json` flag to `analyze.go`.