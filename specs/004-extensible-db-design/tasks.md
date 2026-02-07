# Tasks: Release Readiness, Safety & Extensible Database Design

**Input**: Design documents from `/specs/004-extensible-db-design/`
**Prerequisites**: plan.md, spec.md

## Phase 1: Database Abstraction (Foundational) ✅

- [x] T001 Define `DB` and `Tx` interfaces in `internal/db/interface.go`
- [x] T002 Move PostgreSQL implementation to `internal/db/postgres/client.go` and implement interface
- [x] T003 Implement Backend Registry and DSN detection in `internal/db/registry.go`
- [x] T004 Verify `task ready` still passes after abstraction (regression check)

---

## Phase 2: Engine & Safety (Core Logic) ✅

- [x] T005 Refactor `internal/masker/engine.go` to depend exclusively on `db.DB`
- [x] T006 Implement strict `DryRun` enforcement in engine (ensure zero writes)
- [x] T007 Add `EstimateRowCount` to `DB` interface and engine logic
- [x] T008 Enhance `Engine.Apply` to return a structured `MaskingReport` (tables, estimates, diffs)

---

## Phase 3: CLI Ergonomics & Verification ✅

- [x] T009 Update `cmd/dbmask/apply.go` with `--dry-run` flag and enhanced safety gates
- [x] T010 Update `cmd/dbmask/plan.go` to display table summaries and row estimates
- [x] T011 Update `cmd/dbmask/doctor.go` to verify DSN support and dry-run availability
- [x] T012 Implement `internal/db/postgres.GetDatabaseName` improvement (implicit selection)

---

## Phase 4: Final Verification & Release Readiness ✅

- [x] T013 Update E2E tests to assert zero database modifications in dry-run mode
- [x] T014 Verify adding a mock database implementation requires no changes to engine code
- [x] T015 Final quality check with `task all`