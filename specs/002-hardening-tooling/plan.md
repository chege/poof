# Implementation Plan: hardening-tooling

**Branch**: `002-hardening-tooling` | **Date**: 2026-02-07 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/002-hardening-tooling/spec.md`

## Summary

This phase focuses on introducing `Taskfile.yml` for standardized orchestration, expanding the deterministic faker providers, and hardening the core engine with stricter validation and explicit error handling.

## Technical Context

**Language/Version**: Go 1.25.x
**Primary Dependencies**: 
- `github.com/spf13/cobra`
- `github.com/hashicorp/hcl/v2`
- `github.com/go-faker/faker/v4`
- `github.com/testcontainers/testcontainers-go`
**New Tooling**: Taskfile.dev
**Testing**: `task check` (fmt -> vet -> test)
**Constraints**: No elevation; No new CLI commands; Strict determinism.

## Constitution Check

*Passes Core Principles:*
- Tooling-First: Tasks will guide all autonomous work.
- Determinism: New providers must be seeded per row.
- Hardened: Explicit validation for HCL unknown fields.

## Project Structure

### Documentation (this feature)

```text
specs/002-hardening-tooling/
├── plan.md              # This file
├── spec.md              # Feature specification
├── analyze.md           # Architectural refinements
└── tasks.md             # Task list
```

### Source Code Changes

```text
Taskfile.yml             # New: Taskfile.dev orchestration
internal/config/
  hcl.go                 # Update: Add strict validation for unknown fields
internal/generator/
  all.go                 # Update: Register new faker providers
  faker.go               # Update: Implement new providers
  test_fakers.go         # Update: Add test fakers for new providers
internal/masker/
  engine.go              # Update: Harden error paths and validation
```

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | N/A | N/A |
