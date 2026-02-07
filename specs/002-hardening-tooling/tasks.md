# Tasks: hardening-tooling

**Input**: Design documents from `/specs/002-hardening-tooling/`
**Prerequisites**: plan.md, spec.md

## Phase 1: Setup (Tooling) ✅

- [x] T001 Create `Taskfile.yml` with `build`, `test`, `fmt`, `vet`, `check` tasks
- [x] T002 Implement `init`, `plan`, `apply` tasks in `Taskfile.yml`
- [x] T003 Verify `task check` runs successfully on current master state

---

## Phase 2: Foundational (Hardening) ✅

- [x] T004 Update `internal/config/hcl.go` to use `hcl.Body.Content` for strict validation (fail on unknown fields)
- [x] T005 Refactor `internal/poof/engine.go` to improve error messages and ensure all implicit behaviors are explicit
- [x] T006 Remove any identified dead code or unused abstractions in `internal/`

---

## Phase 3: User Story 2 - Faker Extensions (Priority: P2) ✅

- [x] T007 Implement `username` provider in `internal/generator/faker.go`
- [x] T008 Implement `company_name` provider in `internal/generator/faker.go`
- [x] T009 Implement `phone_number` provider in `internal/generator/faker.go`
- [x] T010 Implement `ipv4_address` provider in `internal/generator/faker.go`
- [x] T011 Implement `short_text` provider in `internal/generator/faker.go`
- [x] T012 Register all new providers in `internal/generator/all.go`
- [x] T013 Add test fakers for all new providers in `internal/generator/test_fakers.go`

---

## Phase 4: Verification & Documentation ✅

- [x] T014 Update E2E tests in `internal/poof/engine_test.go` to cover new providers
- [x] T015 Create a deterministic demo dataset (SQL/HCL) for documentation examples
- [x] T016 Update `README.md` with Taskfile usage and new provider list
- [x] T017 Final verification run with `task check`