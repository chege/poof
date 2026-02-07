# Tasks: Comfortable CLI

**Input**: Design documents from `/specs/003-comfortable-cli/`
**Prerequisites**: plan.md, spec.md

## Phase 1: Setup & UI Infrastructure ✅

- [x] T001 Add UI dependencies: `github.com/fatih/color`, `github.com/mattn/go-isatty`
- [x] T002 Implement TTY detection and symbol-based output in `internal/ui/output.go`
- [x] T003 Add `--no-color` global flag to `cmd/poof/root.go`

---

## Phase 2: Foundational Commands (Non-Destructive) ✅

### User Story 1 - poof doctor (P1)
- [x] T004 Implement readiness checks (config, db, safety, providers) in `internal/poof/doctor.go`
- [x] T005 Implement `poof doctor` command in `cmd/poof/doctor.go`

### User Story 2 - poof init (P2)
- [x] T006 Create annotated HCL template in `internal/config/template.go`
- [x] T007 Implement `poof init` command with `--explain` in `cmd/poof/init.go`

### User Story 3 - poof validate & plan (P3)
- [x] T008 Implement `poof validate` in `cmd/poof/validate.go`
- [x] T009 Extend `internal/poof/engine.go` to support `DryRun` mode (returning sample diffs)
- [x] T010 Implement `poof plan` in `cmd/poof/plan.go` with sample output

---

## Phase 3: Application Ergonomics ✅

### User Story 4 - poof apply (P4)
- [x] T011 Update `cmd/poof/apply.go` to always run validation and plan first
- [x] T012 Implement `--yes` flag to skip plan summary output
- [x] T013 Update `apply` output to use the new UI infrastructure (colors/symbols)

---

## Phase 4: Verification ✅

- [x] T014 Update E2E tests to verify `doctor` and `plan` outputs
- [x] T015 Verify `init` produces a valid config that passes `validate`
- [x] T016 Final verification with `task check`