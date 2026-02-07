# Implementation Plan: Data Producers & Row Sources

**Branch**: `005-data-producers` | **Date**: 2026-02-07 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/005-data-producers/spec.md`

## Summary

This phase introduces a pluggable `Producer` architecture to decouple row selection from the masking engine. The engine will consume a stream of primary keys from a producer, allowing for standard table scans, view-based filtering, and custom SQL queries.

## Technical Context

**Language/Version**: Go 1.25.x
**Primary Dependencies**: 
- `github.com/jackc/pgx/v5`
- `github.com/hashicorp/hcl/v2`
**Architecture**: Interface-driven row selection.
**Constraints**: Compile-time registration; No external data sources; Absolute determinism.

## Constitution Check

*Passes Core Principles:*
- Row Selection Extensibility: Decoupled via `Producer` interface.
- Determinism: Ordering enforced at the producer level.
- Safety: Producers must respect database boundaries and allowlists.

## Project Structure

### Documentation

```text
specs/005-data-producers/
├── plan.md              # This file
├── spec.md              # Feature specification
└── checklists/
    └── requirements.md  # Quality checklist
```

### Source Code Changes

```text
internal/producer/
  interface.go           # New: Producer and Registry definitions
  registry.go            # New: Compile-time producer registration
  table.go               # New: Standard table scan producer
  view.go                # New: View-based producer
  query.go               # New: Custom SQL query producer
internal/config/
  models.go              # Update: Add Source block to Table struct
internal/poof/
  engine.go              # Update: Consume rows from Producer interface
```

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | N/A | N/A |
