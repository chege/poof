# Implementation Plan: poof-cli

**Branch**: `001-poof-cli` | **Date**: 2026-02-07 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-poof-cli/spec.md`

## Summary

Build a Go-based CLI tool `poof` that performs deterministic, declarative data masking for PostgreSQL. The tool will use HCL for configuration, `pgx` for database interaction, and a compile-time registry for extensible generators.

## Technical Context

**Language/Version**: Go 1.25.x
**Primary Dependencies**: 
- `github.com/spf13/cobra` (CLI)
- `github.com/hashicorp/hcl/v2` (Config)
- `github.com/jackc/pgx/v5` (PostgreSQL)
- `github.com/go-faker/faker/v4` (Data Generation)
- `github.com/testcontainers/testcontainers-go` (Testing)
**Storage**: PostgreSQL (Target)
**Testing**: `go test ./...` with Testcontainers-go for E2E.
**Target Platform**: Linux (Debian)
**Project Type**: CLI
**Performance Goals**: Support parallel masking of rows; deterministic output regardless of worker count.
**Constraints**: No runtime plugin loading; No DDL/Deletes; One transaction per table.

## Constitution Check

*Passes Core Principles:*
- Library-First: Logic will be in `internal/` packages.
- CLI Interface: Cobra-based CLI in `cmd/`.
- Test-First: E2E tests with real PostgreSQL containers.
- Determinism: Seeded with MD5 of PK.

## Project Structure

### Documentation (this feature)

```text
specs/001-poof-cli/
├── plan.md              # This file
├── spec.md              # Feature specification
└── checklists/
    └── requirements.md  # Quality checklist
```

### Source Code (repository root)

```text
cmd/
  poof/
    main.go              # Entry point
    root.go              # Root command
    apply.go             # Apply command implementation
internal/
  config/
    hcl.go               # HCL parsing logic
    models.go            # Configuration structs
  generator/
    registry.go          # Generator registration system
    interface.go         # Generator interface definition
    faker.go             # Faker generator implementation
    template.go          # Template generator
    constant.go          # Constant generator
    null.go              # Null generator
    test_fakers.go       # Predictable fakers for testing
  poof/
    engine.go            # Orchestrates the masking process
    worker.go            # Parallel row processing
  db/
    postgresql.go        # PostgreSQL specific operations (pgx)
```

**Structure Decision**: Standard Go project layout with `cmd/` for binaries and `internal/` for private libraries. This prevents external importing of internal logic and keeps the repository clean.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | N/A | N/A |
