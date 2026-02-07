# Implementation Plan: TOML Configuration, Autonomous & Safe Inline Masking

**Branch**: `006-toml-config-safety` | **Date**: 2026-02-07 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/006-toml-config-safety/spec.md`

## Summary

This phase transitions the tool's configuration to TOML, enforces strict safety checks for in-place masking, and ensures full autonomy. We will remove HCL support and implement a robust TOML-based configuration system that targets a single database per file.

## Technical Context

**Language/Version**: Go 1.25.x
**Primary Dependencies**: 
- `github.com/BurntSushi/toml` (for strict TOML parsing)
- `github.com/spf13/cobra`
- `github.com/jackc/pgx/v5`
**Architecture**: Engine depends on `db.DB` interface; Configuration maps directly to TOML structures.
**Constraints**: No elevation; Irreversible mutation; Mandatory Dry-Run.

## Constitution Check

*Passes Core Principles:*
- TOML-Only: Reduces complexity and removes DSL ambiguity.
- Safety-First: Plan and Dry-Run are mandatory gates.
- Autonomous: No prompts or confirmations.

## Project Structure

### Documentation

```text
specs/006-toml-config-safety/
├── plan.md              # This file
├── spec.md              # Feature specification
└── checklists/
    └── requirements.md  # Quality checklist
```

### Source Code Changes

```text
internal/config/
  models.go              # Update: TOML tags and structure
  toml.go                # New: TOML parsing and strict validation
  hcl.go                 # Delete: HCL support removal
cmd/poof/
  root.go                # Update: Use poof.toml by default
  init.go                # Update: Generate TOML template
  validate.go            # Update: Validate TOML
  apply.go               # Update: Enforce DryRun/Safety gates
internal/poof/
  engine.go              # Update: Single DB logic, report refinements
```

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Removal of HCL | Strict compliance with SRS 6.0 | Keeping both would increase maintenance and cognitive load |
