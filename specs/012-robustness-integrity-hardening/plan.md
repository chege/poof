# Implementation Plan: Robustness and Integrity Hardening

**Branch**: `012-robustness-integrity-hardening` | **Date**: Sunday, February 8, 2026 | **Spec**: [specs/012-robustness-integrity-hardening/spec.md](specs/012-robustness-integrity-hardening/spec.md)
**Input**: Feature specification from `specs/012-robustness-integrity-hardening/spec.md`

## Summary

This plan introduces deep validation and safety markers to ensure data integrity. Key components include type-compatibility verification between generators and SQL types, a database-level state marker to detect partial masking runs, in-memory logical uniqueness enforcement, and circular dependency detection during schema analysis.

## Technical Context

**Language/Version**: Go 1.25.7
**Primary Dependencies**:
- `internal/db` (Metadata and State management)
- `internal/engine` (Worker loop enhancements)
- `internal/analyze` (Cycle detection)
**Database**: PostgreSQL 14+

## Constitution Check

- **Safe-by-default**: Pass. State markers prevent accidental use of partially masked data.
- **No Magic**: Pass. All safety checks are explicit and explainable.
- **Deterministic**: Pass. Logical uniqueness retries will use deterministic seeding.

## Project Structure

### Source Code

```text
internal/
├── db/
│   ├── interface.go    # Update: Add GetForeignKeys, GetState, SetState
│   └── postgres/
│       ├── schema.go   # Update: Implement GetForeignKeys
│       └── state.go    # New: Manage _poof_state table
├── engine/
│   ├── engine.go       # Update: Manage job state, handle force_unique
│   └── worker.go       # Update: Enforce logical uniqueness
├── generator/
│   └── interface.go    # Update: Add ExpectedType() to Generator interface
└── analyze/
    ├── analyzer.go     # Update: Warn on circular dependencies
    └── graph.go        # New: Cycle detection logic
```

## Research (Phase 0)

### 1. Type Compatibility Mapping
I will implement a mapping between Go types produced by generators and PostgreSQL types.
- `string` -> `text`, `varchar`, `uuid`, `json`, `jsonb`, `timestamp`
- `int64` -> `integer`, `bigint`, `numeric`
- `bool` -> `boolean`
- `[]byte` -> `bytea`

### 2. Job State Marker (`_poof_state`)
Instead of a complex history, we'll use a single-row table or a specific marker.
**Table Schema**:
```sql
CREATE TABLE IF NOT EXISTS _poof_state (
    id SERIAL PRIMARY KEY,
    job_id UUID,
    status TEXT, -- 'STARTED', 'COMPLETED', 'FAILED'
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### 3. Logical Uniqueness
The `Engine` will maintain a `map[string]map[any]bool` (columnKey -> seenValues).
**Caveat**: Large tables (1M+ rows) will consume significant memory if masking many unique columns. We will add a warning if `force_unique` is used on tables with high estimated row counts.

## Data Model

### Configuration Extensions

```toml
[[tables.columns]]
name = "user_slug"
generator = "short_text"
force_unique = true
```

## Detailed Tasks (Phase 2 Preview)

1.  **Metadata**: Update `Generator` interface to return `ExpectedType()`.
2.  **DB Layer**: Implement `_poof_state` management in `internal/db/postgres/state.go`.
3.  **Validation**: Update `ValidateDatabase` to check type compatibility and state markers.
4.  **Engine**: Implement `force_unique` retry logic in `maskTable`.
5.  **Analyze**: Implement cycle detection for foreign keys.