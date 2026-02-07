# Implementation Plan: Release Readiness, Safety & Extensible Database Design

**Branch**: `004-extensible-db-design` | **Date**: 2026-02-07 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/004-extensible-db-design/spec.md`

## Summary

This phase hardens the `poof` tool by abstracting the database layer behind an interface, enabling future support for other databases while isolating the current PostgreSQL implementation. It also enhances the safety features with a mandatory dry-run capability and a detailed planning output.

## Technical Context

**Language/Version**: Go 1.25.x
**Primary Dependencies**: 
- `github.com/jackc/pgx/v5`
- `github.com/hashicorp/hcl/v2`
**Storage**: PostgreSQL (Target)
**Testing**: `task ready` (includes E2E with Testcontainers)
**Constraints**: No ORMs; Minimal interfaces; Implicit backend selection via DSN.

## Constitution Check

*Passes Core Principles:*
- Library-First: Database abstraction lives in `internal/db`.
- Safety-First: Irreversible mutations require plan/dry-run verification.
- Extensible: Design allows new backends without engine changes.

## Project Structure

### Documentation

```text
specs/004-extensible-db-design/
├── plan.md              # This file
├── spec.md              # Feature specification
└── checklists/
    └── requirements.md  # Quality checklist
```

### Source Code Changes

```text
internal/db/
  interface.go           # New: Minimal database interface
  registry.go            # New: Backend registration and selection
  postgres/              # New: Isolated PostgreSQL package
    client.go            # Move from internal/db/postgresql.go
internal/poof/
  engine.go              # Update: Depend on db.DB interface, add DryRun support
  plan.go                # New: Plan generation logic (estimates, diffs)
cmd/poof/
  apply.go               # Update: Mandatory safety checks, --dry-run flag
  plan.go                # Update: Detailed plan output
  doctor.go              # Update: DSN support and dry-run capability check
```

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Interface Abstraction | Enable future DBs without engine rewrite | Direct implementation is easier but blocks extensibility |
