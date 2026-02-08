# Implementation Plan: Config Pre-flight Validation

**Branch**: `011-config-preflight-validation` | **Date**: Sunday, February 8, 2026 | **Spec**: [specs/011-config-preflight-validation/spec.md](specs/011-config-preflight-validation/spec.md)
**Input**: Feature specification from `specs/011-config-preflight-validation/spec.md`

## Summary

This plan enhances the `poof validate` command from a simple syntax check to a deep semantic validator. It will perform live database checks (connectivity, schema existence) and generator parameter validation (template syntax, provider availability) to ensure a configuration is 100% ready for an `apply` run.

## Technical Context

**Language/Version**: Go 1.25.7
**Primary Dependencies**:
- `internal/config` (Validator logic)
- `internal/db` (Schema verification)
- `internal/generator` (Parameter validation)
**Target Platform**: CLI

## Constitution Check

- **Safe-by-default**: Pass. Semantic validation is read-only and prevents runtime masking errors.
- **Boring, Predictable**: Pass. Enhances existing `validate` command with standard pre-flight patterns.

## Project Structure

### Documentation

```text
specs/011-config-preflight-validation/
├── plan.md              # This file
├── spec.md              # Feature specification
└── checklists/
    └── requirements.md
```

### Source Code

```text
cmd/poof/
├── validate.go         # Update: Add --db-check flag and semantic logic

internal/
├── config/
│   ├── models.go
│   └── validator.go    # New: Deep semantic validation logic
├── db/
│   └── interface.go
└── generator/
    └── registry.go     # Update: Add methods to check provider existence
```

## Research (Phase 0)

### 1. Semantic Check Sequence
Validation will occur in tiers:
1.  **Level 1 (Syntax)**: TOML parsing and struct validation (Existing).
2.  **Level 2 (Static Semantic)**:
    -   Validate all `template` strings parse correctly.
    -   Validate all `faker` providers are registered in the global registry.
3.  **Level 3 (Database Semantic - Optional via flag)**:
    -   Connect to DB.
    -   Verify every `table` in config exists.
    -   Verify every `column` in config exists within its table.
    -   Check for `salt` if DB name suggests production.

### 2. Generator Registry Update
I need a way to check if a provider exists without instantiating it.
**Proposed**: `generator.ProviderExists(locale, name) bool`.

## Data Model

### Configuration Extensions
No changes to `Config` struct required. Validation logic will be purely functional.

## Detailed Tasks

1.  **Infrastructure**: Implement `generator.ProviderExists` and `generator.ValidateTemplate`.
2.  **Logic**: Create `internal/config/validator.go` with `ValidateSemantic(ctx, cfg, dbClient)`.
3.  **CLI**: Update `cmd/poof/validate.go` to add `--db-check` (bool) and `--strict` (bool).
4.  **CLI**: Implement logic to load resources and call semantic validator if flag is set.
5.  **Refactor**: Update `cmd/poof/doctor.go` to use the new shared validation logic.