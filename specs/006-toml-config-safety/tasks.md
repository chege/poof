# Tasks: TOML Configuration, Autonomous & Safe Inline Masking

**Input**: Design documents from `/specs/006-toml-config-safety/`
**Prerequisites**: plan.md, spec.md

## Phase 1: Configuration Migration (TOML) ✅

- [x] T001 Add `github.com/BurntSushi/toml` dependency
- [x] T002 Update `internal/config/models.go` with TOML tags and single-database structure
- [x] T003 Implement `LoadConfig` in `internal/config/toml.go` with strict field checking
- [x] T004 Remove HCL parser (`internal/config/hcl.go`) and HCL dependencies
- [x] T005 Update `internal/config/template.go` to produce a documented `poof.toml`

---

## Phase 2: CLI & Engine Refactoring ✅

- [x] T006 Update `cmd/poof/root.go` to use `poof.toml` as default configuration path
- [x] T007 Refactor `cmd/poof/validate.go` to use the new TOML parser
- [x] T008 Update `cmd/poof/init.go` to support TOML initialization
- [x] T009 Refactor `cmd/poof/apply.go` to ensure `dsn` is read from config, and safety gates are non-interactive
- [x] T010 Refactor `cmd/poof/doctor.go` to verify TOML config and connection readiness

---

## Phase 3: Safety & In-Place Mutation ✅

- [x] T011 Enhance `poof.Engine` to strictly enforce `DryRun` (zero writes verified)
- [x] T012 Update `plan` output to include row estimates and generator details from TOML config
- [x] T013 Ensure all commands fail with non-zero exit codes on any validation or safety failure

---

## Phase 4: Verification & Cleanup ✅

- [x] T014 Update E2E tests to use TOML configuration and verify in-place mutation safety
- [x] T015 Verify zero `UPDATE` statements are issued in `--dry-run` mode
- [x] T016 Remove all remaining HCL-related files and sample configurations
- [x] T017 Final verification with `task ready`