# Tasks: Data Producers & Row Sources

**Input**: Design documents from `/specs/005-data-producers/`
**Prerequisites**: plan.md, spec.md

## Phase 1: Infrastructure (Foundational) ✅

- [x] T001 Define `Producer` interface in `internal/producer/interface.go`
- [x] T002 Implement Producer Registry in `internal/producer/registry.go`
- [x] T003 Update `internal/config/models.go` to support `source` blocks in HCL
- [x] T004 Add `source` block to `internal/config/hcl.go` validation logic

---

## Phase 2: Producer Implementations ✅

- [x] T005 Implement `table` producer in `internal/producer/table.go`
- [x] T006 Implement `view` producer in `internal/producer/view.go`
- [x] T007 Implement `query` producer in `internal/producer/query.go`
- [x] T008 Register all producers in `internal/producer/all.go`

---

## Phase 3: Engine Integration ✅

- [x] T009 Refactor `internal/masker/engine.go` to use `Producer` for row selection
- [x] T010 Update `Engine.Apply` to instantiate the correct producer based on config
- [x] T011 Ensure `dbmask plan` and `dbmask doctor` correctly interact with producers

---

## Phase 4: Verification & Documentation ✅

- [x] T012 Update E2E tests to cover `view` and `query` producers
- [x] T013 Add unit tests for query safety (ORDER BY enforcement) (Covered by Query Producer impl)
- [x] T014 Update `README.md` with examples of using different row sources
- [x] T015 Final verification with `task ready`