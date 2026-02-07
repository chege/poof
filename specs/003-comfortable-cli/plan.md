# Implementation Plan: Comfortable CLI

**Branch**: `003-comfortable-cli` | **Date**: 2026-02-07 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/003-comfortable-cli/spec.md`

## Summary

This phase focuses on the CLI ergonomics and user "comfort". We will implement the full set of opinionated commands (`init`, `validate`, `plan`, `doctor`) and improve the `apply` experience with colored, structured output and a confirmation-free flow.

## Technical Context

**Language/Version**: Go 1.25.x
**Primary Dependencies**: 
- `github.com/spf13/cobra`
- `github.com/fatih/color` (for TTY color support)
- `github.com/mattn/go-isatty` (for TTY detection)
**Storage**: PostgreSQL (pgx)
**Testing**: `go test ./...` with Testcontainers-go.
**Constraints**: No elevation; No prompts; No CI in this phase.

## Constitution Check

*Passes Core Principles:*
- CLI-First: Explicit mapping of internal logic to well-defined commands.
- Trust/Predictability: `doctor` and `plan` provide confidence without modification.
- Determinism: Dry-run `plan` uses the same seeding logic as `apply`.

## Project Structure

### Documentation

```text
specs/003-comfortable-cli/
├── plan.md              # This file
├── spec.md              # Feature specification
└── checklists/
    └── requirements.md  # Quality checklist
```

### Source Code Changes

```text
cmd/dbmask/
  root.go                # Update: Global flags (--no-color)
  init.go                # New: Template generation
  validate.go            # New: Config validation
  plan.go                # New: Dry-run and diff preview
  doctor.go              # New: Readiness checks
  apply.go               # Update: Summary and --yes flag
internal/
  ui/
    output.go            # New: Symbols, colors, and TTY logic
  config/
    template.go          # New: HCL template content
  masker/
    engine.go            # Update: Add DryRun mode for plan
```

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | N/A | N/A |
