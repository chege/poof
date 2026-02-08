# Tasks: Poof Core Enhancements

**Input**: Design documents from `/specs/008-poof-core-v1-enhancements/`
**Prerequisites**: plan.md (required), spec.md (required)

## Phase 1: Setup & Foundations

**Purpose**: Establish the infrastructure for structured exits and multi-env configuration.

- [ ] T001 Create `internal/ui/exit.go` with structured exit code constants (ExitOK, ExitConfigErr, ExitConnErr, etc.)
- [ ] T002 Implement `ui.HandleExit(err)` in `internal/ui/exit.go` to map errors to codes and call `os.Exit`
- [ ] T003 [P] Update `internal/config/models.go` to add `Databases` map and `Locale` field to `Config` struct
- [ ] T004 Update `internal/config/toml.go` to handle backward compatibility for the legacy `database` block

---

## Phase 2: User Story 1 - Structured Exit Codes (Priority: P1)

**Goal**: Ensure all CLI commands return category-specific exit codes.

**Independent Test**: Trigger a config error and verify `echo $?` returns 2.

- [ ] T005 [US1] Update `cmd/poof/root.go` to use `ui.HandleExit` for persistent pre-run/post-run errors
- [ ] T006 [US1] Update `cmd/poof/apply.go` to return specific error types (ConfigError, ConnError)
- [ ] T007 [US1] Update `cmd/poof/plan.go`, `cmd/poof/analyze.go`, and `cmd/poof/doctor.go` with structured exit logic
- [ ] T008 [US1] Add unit tests in `internal/ui/exit_test.go` to verify error-to-code mapping

---

## Phase 3: User Story 2 - Multi-Environment Config (Priority: P1) 🎯 MVP

**Goal**: Support selecting database environments via the `--env` flag.

**Independent Test**: Define a "staging" DB and run `poof apply --env staging`.

- [ ] T009 [US2] Add `--env` (aliased to `-e`) persistent flag to `cmd/poof/root.go`
- [ ] T010 [US2] Implement environment selection logic in `internal/config/toml.go` (prioritize `--env`, then `default` flag, then first block)
- [ ] T011 [US2] Update `apply.go` and `plan.go` to use the selected environment's DSN
- [ ] T012 [US2] Add validation to ensure requested environment exists in config

---

## Phase 4: User Story 3 - Localized Masking (Priority: P2)

**Goal**: Support `locale` setting for faker-generated data.

**Independent Test**: Set `locale = "de_DE"` and verify German names are generated.

- [ ] T013 [US3] Refactor `internal/generator/faker.go` to support locale-keyed provider maps
- [ ] T014 [US3] Add basic data slices for `de_DE`, `fr_FR`, and `es_ES` for common providers (names, cities)
- [ ] T015 [US3] Update `fakerGenerator.Generate` to use the locale from `RowContext` (falling back to `en_US`)
- [ ] T016 [US3] Update `internal/generator/context.go` to include `Locale` in `RowContext`

---

## Phase 5: User Story 4 - Generator Composition (Priority: P2)

**Goal**: Allow templates to invoke other generators (e.g., `{{ faker "email" }}`).

**Independent Test**: Use a template `"PREFIX-{{ counter }}"` and verify output.

- [ ] T017 [US4] Update `internal/generator/template.go` to include a `FuncMap` in the template execution
- [ ] T018 [US4] Implement `faker`, `counter`, and `hash` functions for use within templates
- [ ] T019 [US4] Add unit tests in `internal/generator/template_test.go` for composed outputs
- [ ] T020 [US4] Update documentation/examples for template composition

---

## Phase 6: Polish & Verification

- [ ] T021 Update `demo/poof.toml` to showcase multiple environments and localized masking
- [ ] T022 Run `task ready` to verify all enhancements and existing features
- [ ] T023 Final README.md updates for exit codes and multi-env usage
