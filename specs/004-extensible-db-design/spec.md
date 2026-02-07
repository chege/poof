# Feature Specification: Release Readiness, Safety & Extensible Database Design

**Feature Branch**: `004-extensible-db-design`  
**Created**: 2026-02-07  
**Status**: Draft  
**Input**: User description: "SRS: dbmask — Release Readiness, Safety & Extensible Database Design VERSION: 3.2"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Irreversible Inline Masking with Safety Checks (Priority: P1)

As a Database Administrator, I want to perform in-place masking on a database with absolute confidence that safety checks (allowlist, dry-run) prevent accidental corruption or unintended execution.

**Why this priority**: Inline mutation is irreversible. Safety is the highest priority to prevent data loss or unauthorized masking.

**Independent Test**: Can be tested by attempting to run `dbmask apply` on a database not in the allowlist (should refuse) and by running with `--dry-run` (should perform all logic but commit no changes).

**Acceptance Scenarios**:

1. **Given** a database name not in the `allowlist`, **When** I run `dbmask apply`, **Then** the tool refuses to proceed and exits with an error.
2. **Given** a valid configuration, **When** I run `dbmask apply --dry-run`, **Then** the tool connects, fetches data, computes masks, but performs zero writes to the database.

---

### User Story 2 - Verifiable Masking Plan (Priority: P2)

As a Developer, I want to see a detailed summary of affected tables, columns, and a preview of data changes before I execute an irreversible masking operation.

**Why this priority**: Build trust and allow human verification of the masking logic before it is applied to the data.

**Independent Test**: Run `dbmask plan` and verify the output contains table names, estimated row counts, generator types, and a limited before/after sample.

**Acceptance Scenarios**:

1. **Given** a valid configuration, **When** I run `dbmask plan`, **Then** I see a summary listing each table, its estimated row count, and the generator being used for each masked column.
2. **Given** sample data, **When** I run `dbmask plan`, **Then** I see a limited number of rows (e.g., 5) showing the current value and the generated masked value for verification.

---

### User Story 3 - Extensible Database Backend (Priority: P3)

As a Software Engineer, I want the database interaction logic to be abstracted behind a clean interface, so that supporting new database types in the future does not require refactoring the core masking engine.

**Why this priority**: Ensures long-term maintainability and architectural cleaness. Isolates SQL dialect specific logic.

**Independent Test**: Verify that the core engine code (`internal/masker`) depends only on interfaces defined in `internal/db` and that PostgreSQL-specific code is isolated in its own package.

**Acceptance Scenarios**:

1. **Given** the codebase, **When** I inspect `internal/masker/engine.go`, **Then** I see no direct references to `pgx` or PostgreSQL-specific types.
2. **Given** a new database implementation that fulfills the `Database` interface, **When** registered, **Then** it can be used by the engine without modifications to the engine's core logic.

---

### Edge Cases

- **Unsupported DSN**: The tool must fail fast if the provided DSN format is not recognized or supported by any registered backend.
- **Dry-Run Failure**: If a dry-run capability check fails (e.g., read-only connection when writes would be needed for a real run), the tool must report this in `doctor`.
- **Large Table Row Counts**: Use efficient ways to estimate row counts (e.g., `EXPLAIN` or metadata) rather than `COUNT(*)` where possible to avoid performance hits during planning.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Introduce a minimal database interface (`internal/db/interface.go`) covering schema introspection, row selection (ordered by PK), prepared updates, and transaction handling.
- **FR-002**: Isolate PostgreSQL implementation in `internal/db/postgres`.
- **FR-003**: The core engine MUST depend on the database interface, not concrete implementations.
- **FR-004**: `dbmask apply --dry-run` MUST perform all masking logic (fetch, generate) but execute zero updates/commits.
- **FR-005**: `dbmask plan` MUST display affected tables, columns, estimated row counts, and generator types.
- **FR-006**: `dbmask plan` MUST show a limited (max 5 rows) before → after diff preview.
- **FR-007**: `dbmask doctor` MUST verify DSN support, connectivity, and dry-run capability.
- **FR-008**: Database backend selection MUST be implicit via DSN parsing.

### Key Entities

- **Database Interface**: The abstraction layer defining mandatory operations for any backend.
- **PostgreSQL Backend**: The concrete implementation of the database interface for PostgreSQL.
- **Dry-Run Engine**: A configuration of the engine that suppresses write operations.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: `dbmask apply --dry-run` results in zero `UPDATE` statements being executed on the database server.
- **SC-002**: Core engine logic (`internal/masker`) remains unchanged when a new database backend is added (only wiring/registration changes).
- **SC-003**: `dbmask plan` completes in under 5 seconds for a schema with 10 tables (excluding data fetching).
- **SC-004**: `dbmask doctor` correctly identifies and flags unsupported connection strings.