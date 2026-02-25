# Tasks: poof-cli

**Input**: Design documents from `/specs/001-poof-cli/`
**Prerequisites**: plan.md, spec.md

## Phase 1: Setup (Shared Infrastructure) ✅

**Purpose**: Project initialization and basic structure

- [x] T001 Initialize Cobra boilerplate in `cmd/poof/main.go` and `cmd/poof/root.go`
- [x] T002 Add dependencies to `go.mod`: `toml`, `pgx/v5`, `go-faker/v4`, `testcontainers-go`
- [x] T003 Create directory structure: `internal/config`, `internal/generator`, `internal/poof`, `internal/db`

---

## Phase 2: Foundational (Blocking Prerequisites) ✅

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

- [x] T004 Implement TOML configuration models in `internal/config/models.go`
- [x] T005 Implement TOML parser in `internal/config/toml.go`
- [x] T006 Implement Generator interface and Registry in `internal/generator/registry.go` and `interface.go`
- [x] T007 Implement base generators: `constant`, `null` in `internal/generator/`
- [x] T008 Implement PostgreSQL client wrapper using `pgx` in `internal/db/postgresql.go`
- [x] T009 Create RowContext and seeding logic in `internal/generator/context.go`

---

## Phase 3: User Story 1 - Mask Sensitive Production Data (Priority: P1) ✅

**Goal**: Implement basic masking functionality with TOML config and faker generators.

- [x] T010 Implement `faker` generator in `internal/generator/faker.go`
- [x] T011 Implement `apply` command logic in `cmd/poof/apply.go`
- [x] T012 Implement Safety Check (allowlist vs --force) in `cmd/poof/apply.go`
- [x] T013 Implement Table masking orchestration in `internal/poof/engine.go`
- [x] T014 Add basic error handling for missing columns or tables

---

## Phase 4: User Story 2 - Deterministic Masking (Priority: P2) ✅

**Goal**: Ensure masking is deterministic based on PK and supports parallelism.

- [x] T015 Implement MD5 seeding based on `table_name:pk_value` in `internal/generator/context.go`
- [x] T016 Implement Worker Pool for parallel row processing in `internal/poof/engine.go`
- [x] T017 Ensure `ORDER BY pk` is used in all selection queries in `internal/db/postgresql.go`
- [x] T018 Verify thread-safety of generator registry and faker implementations

---

## Phase 5: User Story 3 - Predictable Testing (Priority: P3) ✅

**Goal**: Implement test fakers and E2E integration tests.

- [x] T019 Implement dedicated test faker providers in `internal/generator/test_fakers.go`
- [x] T020 Setup E2E test suite using `testcontainers-go` in `internal/poof/engine_test.go`
- [x] T021 Write E2E test cases for all generator types and safety behaviors
- [x] T022 Add unit tests for TOML parsing and Registry registration (Covered by E2E)

---

## Phase 6: Polish & Cross-Cutting Concerns ✅

**Purpose**: Final cleanup and documentation.

- [x] T023 Add documentation comments to all exported functions
- [x] T024 Code cleanup and refactoring
- [x] T025 Verify "fail hard" behavior (Addressed via TOML parsing)
- [x] T026 Final build check with `go build ./cmd/poof`