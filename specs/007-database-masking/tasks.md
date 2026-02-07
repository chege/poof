# Tasks: Database Masking CLI (Poof)

**Input**: Design documents from `/specs/007-database-masking/`
**Prerequisites**: plan.md (required), spec.md (required)

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [ ] T001 [P] Create `internal/analyze/` directory for the heuristic engine
- [ ] T002 [P] Create `internal/db/postgres/schema.go` for database introspection helpers

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

- [ ] T003 Enhance `internal/engine/engine.go` to support a non-limited `DryRun` mode (removing the 5-row hardcoded limit)
- [ ] T004 Update `internal/engine/engine.go` to collect and return detailed statistics (updated, retried, failed counts)
- [ ] T005 [P] Implement `hash` generator in `internal/generator/hash.go`
- [ ] T006 [P] Implement `counter` generator in `internal/generator/counter.go`
- [ ] T007 Implement the retry loop skeleton in `internal/engine/worker.go` to catch unique constraint errors

**Checkpoint**: Foundation ready - core engine and missing generators are in place.

---

## Phase 3: User Story 1 - Safe Masking Preview (Priority: P1) 🎯 MVP

**Goal**: Implement `poof plan` to show masking changes without applying them.

**Independent Test**: Run `./poof plan --config demo/poof.toml` and verify the output shows the plan without modifying the DB.

### Implementation for User Story 1

- [ ] T008 Implement `cmd/poof/plan.go` to load config and invoke `engine.Apply` in `DryRun` mode
- [ ] T009 [P] Update `internal/ui/output.go` to render a clean, tabular execution plan
- [ ] T010 Add validation in `plan.go` to ensure all configured tables and columns exist in the target schema

**Checkpoint**: User Story 1 (Plan) is fully functional.

---

## Phase 4: User Story 2 - Deterministic Masking & Apply (Priority: P1)

**Goal**: Ensure `poof apply` is deterministic and handles the `--dry-run` flag.

**Independent Test**: Run `./poof apply --dry-run` twice and verify output values are identical.

### Implementation for User Story 2

- [ ] T011 Update `cmd/poof/apply.go` to respect the `--dry-run` flag and call the engine appropriately
- [ ] T012 Ensure `internal/engine/engine.go` uses the row's Primary Key as a seed modifier for all generators
- [ ] T013 Update `internal/ui/output.go` to display final statistics (rows updated, failed, etc.)

---

## Phase 5: User Story 3 - Unique Constraint Retries (Priority: P2)

**Goal**: Recover from unique violations using deterministic retries.

**Independent Test**: Configure a masking rule on a UNIQUE column and verify retries in the logs/output.

### Implementation for User Story 3

- [ ] T014 Implement PostgreSQL `23505` error detection in `internal/engine/worker.go`
- [ ] T015 Implement the deterministic retry seed logic: `generator.Generate(row, attempt_count)`
- [ ] T016 Add `MaxRetries` support to the engine configuration and default to 10

---

## Phase 6: User Story 4 - Advisory Masking Analysis (Priority: P3)

**Goal**: Implement `poof analyze` to suggest masking candidates.

**Independent Test**: Run `./poof analyze` and verify it suggests common sensitive columns (email, etc.).

### Implementation for User Story 4

- [ ] T017 Implement schema introspection in `internal/db/postgres/schema.go` to fetch all column names
- [ ] T018 Implement regex-based heuristic rules in `internal/analyze/rules.go`
- [ ] T019 Implement `cmd/poof/analyze.go` to print suggested TOML snippets to stdout

---

## Phase 7: Polish & Demo

**Purpose**: Final verification and documentation.

- [ ] T020 Update `demo/poof.toml` and `demo/demo.sql` to include a UNIQUE constraint example
- [ ] T021 Run `task ready` to ensure all tests, linting, and formatting pass
- [ ] T022 [P] Final README.md updates for new commands (`plan`, `analyze`)

---

## Dependencies & Execution Order

1. **Setup & Foundational (Phases 1-2)**: MUST be completed first.
2. **User Story 1 (P1)**: The primary MVP deliverable.
3. **User Story 2 & 3 (P1/P2)**: Can proceed in parallel after US1.
4. **User Story 4 (P3)**: Lowest priority, independent of the masking engine.
