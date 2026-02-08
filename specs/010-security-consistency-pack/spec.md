# Feature Specification: Security & Consistency Pack

**Feature Branch**: `010-security-consistency-pack`  
**Created**: Saturday, February 7, 2026  
**Status**: Draft  
**Input**: User description: "Cross-Table Consistency, Graceful Shutdown, Global Seeding Salt, Incremental Masking, JSON Analysis Export"

## User Scenarios & Testing

### User Story 1 - Relational Integrity via Consistent Seeding (Priority: P1)

As a data architect, I want to ensure that if a user's email appears in multiple tables (e.g., `users` and `orders`), it is masked to the same value in both places, so that I can maintain referential integrity and perform joins on masked data for analysis.

**Why this priority**: Essential for testing complex applications where data is denormalized or relies on consistent foreign keys.

**Independent Test**: Configure two tables with the same column (e.g., `email`) and use `seed_by = "value"`. Verify that identical original values result in identical masked values across both tables.

**Acceptance Scenarios**:

1. **Given** two tables with a common `email` column, **When** masked with `seed_by = "value"`, **Then** the row `bob@example.com` in Table A is masked to the exact same value as `bob@example.com` in Table B.
2. **Given** a standard config (default), **When** masked, **Then** behavior remains `table + pk` seeding (backward compatible).

---

### User Story 2 - Graceful Shutdown (Priority: P1)

As an operator, I want the masking process to handle `Ctrl+C` gracefully by finishing the current batch and rolling back the transaction, so that I don't leave the database in an inconsistent or locked state.

**Why this priority**: Prevents operational headaches and database locks during long-running jobs.

**Independent Test**: Start a long masking job and send `SIGINT`. Verify that the process exits with a clear message and no "half-committed" transaction artifacts remain.

**Acceptance Scenarios**:

1. **Given** a running masking job, **When** `SIGINT` is received, **Then** the engine stops accepting new batches, waits for the current batch to finish (or roll back), and exits with a structured status report.

---

### User Story 3 - Environment Isolation via Global Salt (Priority: P2)

As a security engineer, I want to enforce different masking outcomes for `dev` and `staging` environments even if the source data is identical, so that a leak in one environment doesn't compromise the other.

**Why this priority**: Enhances security depth by preventing rainbow table attacks across environments.

**Independent Test**: Run masking twice on the same data with two different `[safety] salt` values. Verify that the output is deterministic within each run but different between the two runs.

**Acceptance Scenarios**:

1. **Given** two configs with different `salt` values, **When** applied to identical databases, **Then** the resulting masked values for the same PK are strictly different.

---

### User Story 4 - Incremental Masking (Priority: P2)

As a DBA, I want to mask only the rows that have changed since the last run (e.g., `created_at > yesterday`), so that my nightly masking jobs finish in minutes instead of hours.

**Why this priority**: Critical for scaling to large datasets where full re-masking is too slow.

**Independent Test**: Configure a `filter` on a table and verify that only rows matching the condition are touched.

**Acceptance Scenarios**:

1. **Given** a table with old and new rows, **When** a filter `created_at > NOW() - INTERVAL '1 hour'` is applied, **Then** only the new rows are masked.

---

### User Story 5 - Automated Analysis Reporting (Priority: P3)

As a compliance officer, I want to export the schema analysis as a JSON file, so that I can ingest it into my dashboard to track PII coverage across our fleet of databases.

**Why this priority**: Enables automated governance and monitoring.

**Independent Test**: Run `poof analyze --json` and pipe the output to `jq` to verify the structure.

**Acceptance Scenarios**:

1. **Given** a database schema, **When** `poof analyze --json` is executed, **Then** a valid JSON array of suggested rules is printed to stdout.

## Requirements

### Functional Requirements

- **FR-001**: System MUST support `seed_by` configuration per table/column, allowing "value" (content-based) or "global_id" seeding strategies.
- **FR-002**: System MUST intercept `SIGINT` and `SIGTERM` signals to trigger a graceful shutdown sequence.
- **FR-003**: System MUST support a global `salt` string in the `[safety]` configuration block that modifies the seeding algorithm.
- **FR-004**: System MUST allow an optional SQL `filter` clause in the `[tables.source]` configuration to restrict the rows selected for masking.
- **FR-005**: System MUST provide a `--json` flag for the `analyze` command that outputs the analysis results in JSON format.

### Assumptions

- **AS-001**: `seed_by = "value"` requires the original value to be read into memory, which is acceptable for the worker pool architecture.
- **AS-002**: Incremental masking relies on the user providing a valid SQL `WHERE` clause compatible with the target database.

## Success Criteria

### Measurable Outcomes

- **SC-001**: 100% of rows with identical values across tables mask to identical values when `seed_by="value"` is active.
- **SC-002**: `Ctrl+C` results in a clean process exit within 5 seconds (configurable timeout).
- **SC-003**: Changing the global salt changes 100% of the generated masked values.
- **SC-004**: Incremental runs on 1M+ row tables (with 1k changes) complete in <10 seconds.