# Tasks: Config Pre-flight Validation

**Input**: Design documents from `/specs/011-config-preflight-validation/`
**Prerequisites**: plan.md (required), spec.md (required)

## Phase 1: Infrastructure Enhancements

**Purpose**: Add helper methods to the generator and config layers to support deep validation.

- [ ] T001 Update `internal/generator/registry.go` to add `ProviderExists(locale, name) bool` method.
- [ ] T002 Update `internal/generator/template.go` to add `ValidateTemplate(text string) error` helper.
- [ ] T003 [P] Add `internal/config/validator.go` with base `SemanticError` type and interface.

---

## Phase 2: Static Semantic Validation (Level 2)

**Goal**: Validate config parameters that don't require a database connection.

- [ ] T004 Implement `cfg.ValidateStatic()` in `internal/config/validator.go`.
- [ ] T005 [US2] Validate that all `faker` providers used in the config are actually registered.
- [ ] T006 [US2] Validate that all `template` strings in the config are syntactically correct.

---

## Phase 3: Database Semantic Validation (Level 3)

**Goal**: Validate config against the live database schema.

- [ ] T007 Implement `cfg.ValidateDatabase(ctx, dbClient)` in `internal/config/validator.go`.
- [ ] T008 [US1] Verify table existence for all configured tables.
- [ ] T009 [US1] Verify column existence for all configured columns.
- [ ] T010 [US3] Implement production readiness check (salt enforcement for sensitive environments).

---

## Phase 4: CLI Integration

**Goal**: Wire the new validation logic into the `validate` and `doctor` commands.

- [ ] T011 Add `--db-check` and `--strict` flags to `cmd/poof/validate.go`.
- [ ] T012 Update `runValidate()` to support the new validation levels based on flags.
- [ ] T013 [P] Update `cmd/poof/doctor.go` to leverage `ValidateStatic` and `ValidateDatabase`.

---

## Phase 5: Verification

- [ ] T014 Add unit tests in `internal/config/validator_test.go` for various invalid config scenarios.
- [ ] T015 Run `task ready` to ensure no regressions in existing commands.
- [ ] T016 Update README.md with new `validate` usage examples.
