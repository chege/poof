# Feature Specification: TOML Configuration, Autonomous & Safe Inline Masking

**Feature Branch**: `006-toml-config-safety`  
**Created**: 2026-02-07  
**Status**: Draft  
**Input**: User description: "SRS: poof — TOML Configuration, Autonomous & Safe Inline Masking VERSION: 6.0"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Configure Masking with TOML (Priority: P1)

As a Database Administrator, I want to use TOML to configure masking rules for a single database, so that the configuration is explicit, declarative, and easy to audit without the complexity of a DSL.

**Why this priority**: TOML is now the mandatory and exclusive configuration format. This is the primary interface for the tool.

**Independent Test**: Can be tested by creating a `poof.toml` file, running `poof validate`, and ensuring it correctly identifies valid and invalid TOML configurations according to the defined schema.

**Acceptance Scenarios**:

1. **Given** a valid `poof.toml` with `[database]`, `[safety]`, and `[[tables]]` sections, **When** I run `poof validate`, **Then** the tool reports the configuration is valid.
2. **Given** a configuration file with an unknown attribute (e.g. `[unknown] section`), **When** I run `poof validate`, **Then** the tool fails with a clear validation error.

---

### User Story 2 - Safe In-Place Masking with Dry-Run (Priority: P2)

As a DevOps Engineer, I want to execute a dry-run of the masking process before applying changes irreversibly, so that I can verify the masking logic and connection parameters without risk to data.

**Why this priority**: Safety is a core principle. Mandatory dry-run capability provides the necessary guardrail for irreversible operations.

**Independent Test**: Run `poof apply --dry-run` against a test database and verify that while rows are fetched and values generated, zero `UPDATE` statements are committed.

**Acceptance Scenarios**:

1. **Given** a valid TOML configuration, **When** I run `poof apply --dry-run`, **Then** the tool connects to the DB, performs plan calculation, and generates masked values, but exits with a message indicating no changes were written.
2. **Given** a dry-run execution, **When** I inspect the database, **Then** all original data remains unchanged.

---

### User Story 3 - Visualizing the Masking Plan (Priority: P3)

As a Security Auditor, I want to see a clear plan of which tables and columns will be masked, including sample before/after values, so that I can approve the masking strategy.

**Why this priority**: Visibility precedes mutation. A clear plan builds trust and allows for human-in-the-loop verification in sensitive environments.

**Independent Test**: Run `poof plan` and verify the output contains table names, column names, generator types, and a limited preview of data changes.

**Acceptance Scenarios**:

1. **Given** a configured table `users`, **When** I run `poof plan`, **Then** I see an estimate of row counts and a preview showing `email: secret@company.com -> user_123@example.org`.

---

### Edge Cases

- **Multiple Database DSNs**: The system must fail if a user tries to target more than one database per configuration file.
- **Malformed DSN**: `poof doctor` and `poof apply` must fail fast with a clear error if the DSN format is unsupported or invalid.
- **Non-PostgreSQL DSN**: Since only PostgreSQL is supported in this phase, any other DSN scheme must be rejected.
- **Missing Required Fields**: Validation must fail if `dsn`, `allowed_db_names`, or table `pk` are missing.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST use TOML as the sole configuration format.
- **FR-002**: System MUST fail validation on unknown or missing fields in the TOML configuration.
- **FR-003**: System MUST support exactly one database per configuration file.
- **FR-004**: System MUST provide `poof apply --dry-run` which performs zero database writes.
- **FR-005**: `poof plan` MUST show affected tables, columns, row estimates, and sample diffs (max 5 rows).
- **FR-006**: `poof apply` MUST run full validation, planning, and safety checks before mutation.
- **FR-007**: `poof doctor` MUST perform non-destructive pre-flight checks including dry-run capability verification.
- **FR-008**: System MUST support fully autonomous (YOLO) mode with no interactive prompts.
- **FR-009**: System MUST NOT request elevated permissions (sudo).

### Key Entities

- **Configuration (TOML)**: The single source of truth for database DSN, safety allowlists, and masking rules.
- **Plan**: A deterministic, reviewable summary of intended operations.
- **Database Interface**: The internal abstraction layer that isolates database-specific logic.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of configuration operations use TOML; HCL support is removed.
- **SC-002**: `poof apply --dry-run` performs 0 commits to the database.
- **SC-003**: All CLI commands complete without requiring user input or `sudo`.
- **SC-004**: Adding a new database backend later requires zero changes to the core engine logic.
- **SC-005**: `poof doctor` correctly identifies a database not in the `allowed_db_names` list.