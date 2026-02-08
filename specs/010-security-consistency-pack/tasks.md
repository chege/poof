# Tasks: Security & Consistency Pack

**Input**: Design documents from `/specs/010-security-consistency-pack/`
**Prerequisites**: plan.md (required), spec.md (required)

## Phase 1: Foundations & Models

**Purpose**: Update data structures to support new configuration options.

- [ ] T001 Update `internal/config/models.go` to add `Salt` to `Safety`, `SeedBy` to `Column`, and `Filter` to `Source`.
- [ ] T002 Update `internal/generator/interface.go` to add `OriginalValue` to `RowContext`.
- [ ] T003 Update `internal/generator/context.go` to accept `salt` and `originalValue` in `NewRowContext`.

---

## Phase 2: Security & Consistency (Seeding)

**Goal**: Implement global salt and cross-table consistency.

- [ ] T004 Refactor `NewRowContext` in `internal/generator/context.go` to use the new seeding algorithm: `MD5(salt + ":" + identifier)`.
- [ ] T005 Update `internal/engine/engine.go` worker loop to pass the current column's original value to the generator.
- [ ] T006 Update `internal/engine/engine.go` to respect the `SeedBy` configuration when creating `RowContext`.

---

## Phase 3: Performance (Incremental Masking)

**Goal**: Support SQL filters for row selection.

- [ ] T007 Update `internal/db/interface.go` to add an optional `filter` parameter to `FetchRows`.
- [ ] T008 Update `internal/db/postgres/client.go` to implement the `filter` SQL logic in `FetchRows`.
- [ ] T009 Update `internal/producer/table.go` to pass the `Source.Filter` from configuration down to the DB layer.

---

## Phase 4: Robustness (Graceful Shutdown)

**Goal**: Handle process signals cleanly.

- [ ] T010 Update `cmd/poof/root.go` to use `signal.NotifyContext` for `SIGINT` and `SIGTERM`.
- [ ] T011 Ensure `cmd/poof/apply.go` and `plan.go` pass the signal-aware context to `LoadResources`.
- [ ] T012 Update `internal/engine/engine.go` to explicitly log when a shutdown is initiated via context cancellation.

---

## Phase 5: Automation (JSON Analysis)

**Goal**: Support JSON output for the `analyze` command.

- [ ] T013 Add `--json` flag to `cmd/poof/analyze.go`.
- [ ] T014 Implement JSON serialization for `Suggestion` objects in `cmd/poof/analyze.go`.

---

## Phase 6: Polish & Verification

- [ ] T015 Update `demo/poof.toml` with a sample `salt` and `filter`.
- [ ] T016 Add integration test in `internal/engine/engine_test.go` verifying that `seed_by="value"` works across tables.
- [ ] T017 Run `task ready` to verify all enhancements.
