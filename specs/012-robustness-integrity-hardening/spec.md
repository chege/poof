# Feature Specification: Robustness and Integrity Hardening

**Feature Branch**: `012-robustness-integrity-hardening`  
**Created**: Sunday, February 8, 2026  
**Status**: Draft  
**Input**: User description: "Type Mismatch Validation, Global Transaction/Ghost Row Protection, Logical Uniqueness Enforcement, Circular Dependency Detection"

## User Scenarios & Testing

### User Story 1 - Type-Safe Configuration (Priority: P1)

As a developer, I want the validation tool to tell me if I'm trying to use a generator that produces data incompatible with my database column type (e.g., using 'email' on an 'integer' column), so that I don't discover these errors during a production masking run.

**Why this priority**: Prevents high-cost runtime failures and ensures the configuration matches the physical schema reality.

**Independent Test**: Can be tested by configuring a string generator for a numeric column and running `poof validate --db-check`.

**Acceptance Scenarios**:

1. **Given** a column of type `INTEGER`, **When** configured with a `faker` provider that returns `string`, **Then** `poof validate` reports a "Type Mismatch" error.
2. **Given** a column of type `TEXT`, **When** configured with a `counter` (integer), **Then** validation passes (as integers can be cast to text).

---

### User Story 2 - Ghost Row / Partial Success Protection (Priority: P1)

As an auditor, I want to be certain that a database has been completely masked before allowing developers to access it, so that we don't accidentally leak PII from tables that failed during the masking process.

**Why this priority**: Essential for compliance and data privacy. A "half-masked" database is a security breach.

**Independent Test**: Start a multi-table masking job and kill the process. Verify that the database is explicitly marked as "Incomplete" or "Dirty".

**Acceptance Scenarios**:

1. **Given** a masking job starts, **When** it begins, **Then** it creates a marker (e.g., a metadata table `poof_metadata`) indicating `STATUS=MASKING_IN_PROGRESS`.
2. **Given** a job finishes successfully, **When** it completes, **Then** the marker is updated to `STATUS=COMPLETE`.
3. **Given** an incomplete job, **When** a user tries to run `poof doctor`, **Then** it warns that the database is in an inconsistent state.

---

### User Story 3 - Logical Uniqueness for App Slugs (Priority: P2)

As a data engineer, I want to enforce uniqueness on columns like `user_slug` even if they don't have a `UNIQUE` constraint in the database, so that the application logic remains functional after masking.

**Why this priority**: Application stability often depends on unique values that aren't strictly enforced at the database layer (e.g., for performance or legacy reasons).

**Independent Test**: Enable `force_unique` on a column without a constraint and mask 10,000 rows. Verify that no duplicate values exist in the resulting data.

**Acceptance Scenarios**:

1. **Given** a column with `force_unique = true`, **When** the engine generates a value that has already been used in the current run, **Then** it triggers a deterministic retry until a unique value is found.

---

### User Story 4 - Circular Dependency Detection (Priority: P3)

As a DBA, I want to be warned if my schema contains circular dependencies (e.g., Table A points to B, and B points to A), so that I can carefully review the masking order and avoid orphaned records.

**Why this priority**: Helps identify complex relational structures that might require specialized masking logic or specific table ordering.

**Independent Test**: Run `poof analyze` on a schema with a 1:1 circular relationship and verify the warning output.

**Acceptance Scenarios**:

1. **Given** two tables with mutual foreign keys, **When** `poof analyze` is run, **Then** it prints a warning identifying the circular relationship.

---

### Edge Cases

- **Composite PKs**: How does type validation handle columns that are part of a composite primary key?
- **Global Transaction Size**: If using a single transaction for ghost row protection, will the database run out of WAL space for 100M+ row updates?
- **Memory Limit for Uniqueness**: How many unique values can the CLI track in memory before it impacts system performance?

## Requirements

### Functional Requirements

- **FR-001**: System MUST validate that generator output types are compatible with PostgreSQL column data types during semantic validation.
- **FR-002**: System MUST maintain a job state marker in the target database to prevent the use of partially masked data.
- **FR-003**: System MUST support a `force_unique` flag per column to ensure logical uniqueness within a single job run.
- **FR-004**: System MUST detect and report circular Foreign Key dependencies during schema analysis.
- **FR-005**: System MUST provide a mechanism to clean up or "acknowledge" a dirty masking state via the CLI.

### Assumptions

- **AS-001**: Type compatibility mapping will follow standard PostgreSQL implicit casting rules.
- **AS-002**: `force_unique` will be limited by available system memory for tracking used values.

## Success Criteria

### Measurable Outcomes

- **SC-001**: 100% of incompatible type configurations are identified during `poof validate --db-check`.
- **SC-002**: A database that fails mid-masking is correctly identified as "Dirty" by the `doctor` command in 100% of test cases.
- **SC-003**: The `force_unique` implementation ensures zero collisions for up to 1 million unique values per column.
- **SC-004**: Circular dependency detection adds less than 500ms to the `analyze` command execution time.