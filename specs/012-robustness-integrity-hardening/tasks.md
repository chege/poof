# Tasks: Robustness and Integrity Hardening

**Input**: Design documents from `/specs/012-robustness-integrity-hardening/`
**Prerequisites**: plan.md (required), spec.md (required)

## Phase 1: Foundations & Metadata

**Purpose**: Update core interfaces to support deep validation.

- [ ] T001 Update `internal/generator/interface.go` to add `ExpectedType() string` to the `Generator` interface.
- [ ] T002 Implement `ExpectedType()` for all built-in generators (`faker`, `template`, `constant`, `null`, `hash`, `counter`).
- [ ] T003 [P] Add `GetForeignKeys(ctx, tableName) ([]ForeignKey, error)` to `db.DB` interface.
- [ ] T004 Implement `GetForeignKeys` in `internal/db/postgres/schema.go`.

---

## Phase 2: Job State Tracking (Ghost Row Protection)

**Goal**: Implement a database marker to track masking job status.

- [ ] T005 [P] Create `internal/db/postgres/state.go` to manage the `_poof_state` table.
- [ ] T006 Implement `SetJobState(ctx, status)` and `GetJobState(ctx)` in the DB layer.
- [ ] T007 Update `internal/engine/engine.go`: `Apply` method to set state to `STARTED` at the beginning and `COMPLETED` at the end.
- [ ] T008 Update `internal/ui/doctor.go` (or shared validation) to check the state marker and warn if `STARTED` or `FAILED` is found.

---

## Phase 3: Semantic Validation Improvements

**Goal**: Catch type mismatches before they cause runtime errors.

- [ ] T009 [US1] Implement type compatibility checker in `internal/config/validator.go`.
- [ ] T010 Update `ValidateDatabase` to use the type checker and report incompatibilities.
- [ ] T011 [US4] Implement circular dependency detection in `internal/analyze/analyzer.go`.

---

## Phase 4: Logical Uniqueness

**Goal**: Enforce uniqueness even for columns without database constraints.

- [ ] T012 Update `internal/config/models.go` to add `ForceUnique bool` to the `Column` struct.
- [ ] T013 Implement in-memory value tracking in `internal/engine/engine.go` for columns with `ForceUnique`.
- [ ] T014 [US3] Implement retry logic in the worker loop when a logical collision occurs.

---

## Phase 5: Polish & Integration

- [ ] T015 Add unit tests for type compatibility and cycle detection.
- [ ] T016 Run `task ready` to verify all enhancements.
- [ ] T017 Update documentation to explain `force_unique` and the safety marker.
