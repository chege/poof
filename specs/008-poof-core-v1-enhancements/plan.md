# Implementation Plan: Poof Core Enhancements

**Branch**: `008-poof-core-v1-enhancements` | **Date**: Saturday, February 7, 2026 | **Spec**: [specs/008-poof-core-v1-enhancements/spec.md](specs/008-poof-core-v1-enhancements/spec.md)
**Input**: Feature specification from `specs/008-poof-core-v1-enhancements/spec.md`

## Summary

This plan covers the implementation of four high-value enhancements to the Poof core: structured exit codes for automation, multi-environment database support in configuration, localized data masking via faker providers, and generator composition using templates.

## Technical Context

**Language/Version**: Go 1.25.7
**Primary Dependencies**:
- `github.com/spf13/cobra` (CLI)
- `github.com/BurntSushi/toml` (Config)
- `text/template` (Composition)
**Testing**: `stretchr/testify` (Unit)
**Target Platform**: Linux, macOS, Windows (CLI)

## Constitution Check

- **Deterministic**: Pass. Seeding will be maintained for all localized and composed generators.
- **Safe-by-default**: Pass. Allowlist checking will be applied to all environment-specific database blocks.
- **Boring, Predictable**: Pass. Exit codes follow standard POSIX conventions where applicable.

## Project Structure

### Documentation (this feature)

```text
specs/008-poof-core-v1-enhancements/
├── plan.md              # This file
├── spec.md              # Feature specification
└── checklists/          # Quality checklists
    └── requirements.md
```

### Source Code (repository root)

```text
cmd/poof/
├── root.go             # Update: Add --env flag, centralize exit logic
├── apply.go            # Update: Use environment selection, return structured exit codes
└── ...

internal/
├── config/
│   ├── models.go       # Update: Add Databases map, Locale fields
│   └── toml.go         # Update: Handle legacy 'database' vs 'databases'
├── generator/
│   ├── faker.go        # Update: Add locale support to providers
│   ├── template.go     # Update: Add composition context
│   └── context.go
└── ui/
    ├── exit.go         # New: Define structured exit codes
    └── output.go
```

## Research (Phase 0)

### 1. Structured Exit Codes
**Proposed Codes**:
- `0`: Success
- `1`: General Error (fallback)
- `2`: Configuration Error (invalid TOML, missing fields)
- `3`: Connection Error (database unreachable)
- `4`: Partial Masking Failure (job finished but some rows failed)
- `5`: Safety Error (unauthorized database name)

### 2. Multi-Environment Config
**TOML Structure**:
```toml
[databases.local]
dsn = "postgres://..."

[databases.staging]
dsn = "postgres://..."
default = true
```
**Backward Compatibility**: If `[database]` (singular) exists, it is treated as the default environment named "default".

### 3. Locale Support
**Implementation**:
Modify `fakerProviders` map to be nested: `map[string]map[string]FakerProvider` (locale -> providerName -> provider).
If a requested locale is missing a provider, fall back to `en_US`.

### 4. Generator Composition
**Implementation**:
The `template` generator's context will be enhanced with a `FuncMap` allowing it to call other generators by name.
Example: `{{ faker "email" }}` or `{{ counter }}`.

## Data Model (Phase 1)

### Configuration Extensions (`internal/config/models.go`)
```go
type Config struct {
    Database  *Database             `toml:"database"`  // Deprecated
    Databases map[string]Database    `toml:"databases"`
    Locale    string                `toml:"locale"`
    Tables    []Table               `toml:"tables" ...`
}

type Database struct {
    DSN     string `toml:"dsn"`
    Default bool   `toml:"default"`
}
```

## Detailed Tasks (Phase 2 Preview)

1.  **UI**: Implement `internal/ui/exit.go` and update `cmd/poof` to use it.
2.  **Config**: Update `internal/config/models.go` and `internal/config/toml.go` for multi-env support.
3.  **Generator**: Update `faker.go` for localization and `template.go` for composition.
4.  **CLI**: Wire up the `--env` flag and update `apply.go` logic.