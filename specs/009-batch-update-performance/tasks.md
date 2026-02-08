# Tasks: Batch Update Performance

**Input**: Design documents from `/specs/009-batch-update-performance/`
**Prerequisites**: plan.md (required), spec.md (required)

## Phase 1: Setup & Foundational

**Purpose**: Update configuration and prepare the engine for batching.

- [ ] T001 Update `internal/config/models.go` to add `BatchSize` field to `Config` struct.
- [ ] T002 Update `internal/config/toml.go` to set a default `BatchSize` of 500 in `LoadConfig`.
- [ ] T003 Update `internal/engine/engine.go` to include `BatchSize` in the `Engine` struct and `NewEngine` constructor.

---

## Phase 2: User Story 1 - High-Performance Batching (Priority: P1) 🎯 MVP

**Goal**: Implement the core batch update logic using PostgreSQL bulk syntax.

**Independent Test**: Mask a table with 10,000 rows and verify performance gain.

- [ ] T004 [P] Create `internal/engine/batch.go` and implement `buildBatchUpdateQuery` for `UPDATE ... FROM (VALUES ...)`.
- [ ] T005 Refactor `internal/engine/engine.go`: `writeResults` to use a buffer (`[]rowData`) for batching.
- [ ] T006 Implement `internal/engine/engine.go`: `applyBatch` method to execute the bulk SQL within the current transaction.
- [ ] T007 Add benchmark test in `internal/engine/engine_test.go` to measure single vs batch performance.

**Checkpoint**: User Story 1 is functional - masking uses batch updates by default.

---

## Phase 3: User Story 2 - Robust Fallback Handling (Priority: P2)

**Goal**: Recover from batch-level failures by retrying rows individually.

**Independent Test**: Force a unique violation in a batch and verify that other rows are still updated.

- [ ] T008 [US2] Update `applyBatch` to return an error on execution failure.
- [ ] T009 [US2] Implement fallback logic in `writeResults`: on batch failure, iterate through the buffer and call `retryUpdate` for each row.
- [ ] T010 [US2] Add integration test in `internal/engine/engine_test.go` simulating a `UNIQUE` constraint violation inside a large batch.

---

## Phase 4: Polish & Integration

- [ ] T011 Update `demo/poof.toml` to include a sample `batch_size = 100` setting.
- [ ] T012 Run `task ready` to ensure all tests, linting, and formatting pass with the new logic.
- [ ] T013 [P] Update `README.md` to mention the performance boost and configurable batching.

---

## Dependencies & Execution Order

1. **Foundational (Phase 1)**: Must be completed first to provide the `BatchSize` configuration.
2. **User Story 1 (Phase 2)**: Core performance implementation.
3. **User Story 2 (Phase 3)**: Safety and robustness enhancements (depends on US1).
4. **Polish (Phase 4)**: Documentation and final verification.
