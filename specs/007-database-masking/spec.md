# Feature Specification: Database Masking CLI (Poof)

**Feature Branch**: `007-database-masking`  
**Created**: Saturday, February 7, 2026  
**Status**: Draft  
**Input**: User description: "# SRS — Poof (Database Masking CLI) ## 1. Purpose Poof is a safe, opinionated CLI for in-place anonymization (masking) of databases. Primary goal: allow developers to mask real data without breaking constraints, without surprises, and without requiring elevated permissions. Target users: - Developers - Operators - Test / staging environments - Local copies of production data Poof is not a data discovery tool, not an ETL, and not an AI auto-masker. --- ## 2. Non-Goals (Hard Constraints) The system MUST NOT: - Auto-mask data without explicit config - Modify schema (no ALTER, DROP, DISABLE CONSTRAINTS) - Require elevated permissions (SUPERUSER, owner, etc.) - Ask for interactive confirmation mid-run - Perform cloud execution - Perform AI-driven masking decisions - Modify .git/ contents - Execute commands requiring sudo or admin privileges --- ## 3. Supported Databases - PostgreSQL (v1 scope) Design MUST be extensible to support additional databases via code extension (recompile acceptable). --- ## 4. Execution Modes ### 4.1 Plan - Reads schema and config - Produces an execution plan - Does NOT modify data Command: poof plan ⸻ 4.2 Apply • Executes masking based on plan • Supports transactional safety where possible • Explicitly requires invocation Command: poof apply ⸻ 4.3 Dry-Run (Mandatory) • Simulates apply • Generates masked values • Performs no writes Command: poof apply --dry-run Dry-run output MUST match apply output deterministically. ⸻ 4.4 Analyze (Advisory Only) • Read-only inspection • Suggests candidate columns for masking • Uses deterministic, explainable rules Hard rules: • MUST NOT modify DB • MUST NOT generate config • MUST NOT auto-enable masking Command: poof analyze ⸻ 5. Configuration 5.1 Format • TOML (single supported format) • Explicit, readable, version-controllable 5.2 Structure • One or more database blocks • Explicit table and column mapping • No implicit defaults that modify data ⸻ 6. Masking Generators 6.1 Core Generators • Email • Name • Static value • Template value • Hash-based (MD5 / deterministic) • Counter-based 6.2 Generator Properties Generators MUST: • Be deterministic • Support seeded output • Be side-effect free ⸻ 7. Uniqueness & Retry Handling (Critical Requirement) 7.1 Problem Generated masked values may violate UNIQUE constraints. 7.2 Required Behavior • Uniqueness violations are recoverable • Retry is row-scoped, not job-scoped • Retries are bounded • Job MUST continue on per-row failure 7.3 Retry Semantics • Retry attempts generate alternative deterministic values • Retry strategy is generator-controlled • Default retries enabled automatically for UNIQUE columns 7.4 Failure Handling • Exhausted retries mark the row as FAILED • FAILED rows do NOT abort the job • All failures MUST be reported ⸻ 8. Safety Guarantees Poof MUST guarantee: • No schema modification • No constraint disabling • No silent partial failure • Deterministic output across runs • Explicit execution steps • Clear reporting ⸻ 9. Output & Reporting Required Outputs • Rows updated • Rows retried • Rows failed • Failure reasons Exit code: • 0 → success • non-zero → failures occurred (configurable strictness) ⸻ 10. Developer Experience 10.1 Demo Mode Repo MUST include: • Local demo database • Sample config • One-command demo (Taskfile) 10.2 Tooling Project MUST include: • Taskfile.yml • staticcheck • go vet • formatting and lint tasks ⸻ 11. Architecture Principles • Explicit over magic • Deterministic over random • Safe defaults • Opinionated UX • Boring, predictable behavior ⸻ 12. Summary Poof is a trust-first masking tool. It optimizes for: • Safety • Determinism • Predictability • Ease of use All automation MUST be explainable. All destructive actions MUST be explicit."

## User Scenarios & Testing

### User Story 1 - Safe Masking Preview (Priority: P1)

As a developer, I want to see exactly which columns and tables will be masked before any changes are made to the database, so I can verify the configuration.

**Why this priority**: Fundamental safety requirement. Users must trust the tool before applying changes.

**Independent Test**: Can be tested by running `poof plan` with a valid config and verifying the printed execution plan without any database modifications.

**Acceptance Scenarios**:

1. **Given** a database with data and a valid TOML config, **When** `poof plan` is executed, **Then** a detailed list of tables and columns to be masked is displayed, and no data is changed.
2. **Given** an invalid config (e.g., missing table), **When** `poof plan` is executed, **Then** a clear error message is displayed identifying the configuration issue.

---

### User Story 2 - Deterministic Data Masking (Priority: P1)

As a developer, I want the same input data to result in the same masked data every time I run the tool, so that test cases and application behavior remain predictable.

**Why this priority**: Essential for reliable testing and staging environments.

**Independent Test**: Can be tested by running `poof apply --dry-run` twice on the same dataset and comparing the generated values.

**Acceptance Scenarios**:

1. **Given** a specific row in a database, **When** `poof apply` (or dry-run) is executed multiple times, **Then** the masked values for that row are identical in every run.
2. **Given** two identical database instances, **When** `poof apply` is executed on both with the same config, **Then** the resulting masked data is identical.

---

### User Story 3 - Handling Unique Constraint Violations (Priority: P2)

As an operator, I want the masking process to automatically handle cases where a generated value violates a UNIQUE constraint, so that the masking job doesn't fail prematurely.

**Why this priority**: Prevents common masking failures in real-world schemas with constraints.

**Independent Test**: Can be tested by configuring a generator for a column with a UNIQUE constraint and forcing a collision, then verifying the retry logic.

**Acceptance Scenarios**:

1. **Given** a UNIQUE constraint on a column, **When** a generated value collides with an existing one, **Then** the system automatically retries with a new deterministic value.
2. **Given** a collision that persists after the maximum number of retries, **When** the limit is reached, **Then** the row is marked as FAILED, the reason is logged, and the job continues.

---

### User Story 4 - Advisory Masking Analysis (Priority: P3)

As a new user, I want the tool to suggest which columns might contain sensitive data, so that I have a starting point for my masking configuration.

**Why this priority**: Improves discovery and ease of use for new projects.

**Independent Test**: Can be tested by running `poof analyze` on a database with typical sensitive column names (e.g., 'email', 'ssn') and verifying the suggestions.

**Acceptance Scenarios**:

1. **Given** a database with columns like 'email' or 'full_name', **When** `poof analyze` is executed, **Then** it lists these columns as candidates for masking with explainable rules.
2. **Given** `poof analyze` is run, **When** it finishes, **Then** no database data has been modified and no configuration has been automatically generated.

---

### Edge Cases

- **Large Tables**: How does the system handle tables with millions of rows without exhausting memory?
- **Network Interruptions**: How does the system recover if the database connection is lost mid-masking?
- **Mixed Case/Quoted Identifiers**: Does the TOML config correctly handle PostgreSQL identifiers that require double quotes?

## Requirements

### Functional Requirements

- **FR-001**: System MUST support PostgreSQL as the target database.
- **FR-002**: System MUST use TOML for configuration.
- **FR-003**: System MUST provide four execution modes: `plan`, `apply`, `apply --dry-run`, and `analyze`.
- **FR-004**: System MUST NOT require superuser or elevated permissions to perform masking.
- **FR-005**: Generators MUST be deterministic and side-effect free.
- **FR-006**: System MUST support core generators: Email, Name, Static value, Template value, Hash-based, and Counter-based.
- **FR-007**: System MUST automatically retry value generation when a UNIQUE constraint violation occurs.
- **FR-008**: Retries for UNIQUE violations MUST be bounded (max attempts specified or defaulted).
- **FR-009**: System MUST report the count of updated, retried, and failed rows upon completion.
- **FR-010**: System MUST NOT modify the database schema or disable constraints.
- **FR-011**: System MUST provide a "demo mode" with a local DB and sample config for quick evaluation.

### Assumptions

- **AS-001**: The tool will target PostgreSQL 14+ for initial compatibility.
- **AS-002**: A default retry limit of 10 attempts will be used for unique violations unless configured otherwise.
- **AS-003**: Determinism is achieved using a seed derived from the table name and the row's primary key.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Users can run the full demo (setup + plan + apply) in under 5 minutes.
- **SC-002**: Masking a table with 100,000 rows takes less than 60 seconds on standard developer hardware.
- **SC-003**: 100% of masking runs are deterministic; identical inputs produce identical outputs across different environments.
- **SC-004**: Zero schema modifications occur during any masking operation.
- **SC-005**: All unique constraint violations are either resolved via retries or reported as individual row failures without aborting the job.