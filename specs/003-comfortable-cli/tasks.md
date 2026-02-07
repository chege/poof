# Tasks: Comfortable CLI

**Input**: Design documents from `/specs/003-comfortable-cli/`
**Prerequisites**: plan.md, spec.md

## Phase 1: Setup & UI Infrastructure ✅

- [x] T001 Add UI dependencies: `github.com/fatih/color`, `github.com/mattn/go-isatty`
- [x] T002 Implement TTY detection and symbol-based output in `internal/ui/output.go`
- [x] T003 Add `--no-color` global flag to `cmd/dbmask/root.go`

---

## Phase 2: Foundational Commands (Non-Destructive) ✅

### User Story 1 - dbmask doctor (P1)
- [x] T004 Implement readiness checks (config, db, safety, providers) in `internal/masker/doctor.go`
- [x] T005 Implement `dbmask doctor` command in `cmd/dbmask/doctor.go`

### User Story 2 - dbmask init (P2)
- [x] T006 Create annotated HCL template in `internal/config/template.go`
- [x] T007 Implement `dbmask init` command with `--explain` in `cmd/dbmask/init.go`

### User Story 3 - dbmask validate & plan (P3)
- [x] T008 Implement `dbmask validate` in `cmd/dbmask/validate.go`
- [x] T009 Extend `internal/masker/engine.go` to support `DryRun` mode (returning sample diffs)
- [x] T010 Implement `dbmask plan` in `cmd/dbmask/plan.go` with sample output

---

## Phase 3: Application Ergonomics ✅

### User Story 4 - dbmask apply (P4)
- [x] T011 Update `cmd/dbmask/apply.go` to always run validation and plan first
- [x] T012 Implement `--yes` flag to skip plan summary output
- [x] T013 Update `apply` output to use the new UI infrastructure (colors/symbols)

---

## Phase 4: Verification ✅

- [x] T014 Update E2E tests to verify `doctor` and `plan` outputs
- [x] T015 Verify `init` produces a valid config that passes `validate`
- [x] T016 Final verification with `task check`