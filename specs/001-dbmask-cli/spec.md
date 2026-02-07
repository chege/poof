# Feature Specification: dbmask — PostgreSQL Data Masking CLI

**Feature Branch**: `001-dbmask-cli`  
**Created**: 2026-02-07  
**Status**: Draft  
**Input**: User description: "SRS: dbmask — PostgreSQL Data Masking CLI (Go, HCL, Cobra) VERSION: 1.2"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Mask Sensitive Production Data for Development (Priority: P1)

As a Database Administrator, I want to mask sensitive production data (like PII) in a PostgreSQL database using a declarative configuration, so that developers can work with realistic but safe data.

**Why this priority**: This is the core purpose of the tool. Providing a safe way to mask data is the primary value proposition.

**Independent Test**: Can be fully tested by providing an HCL configuration and a target database with sensitive data, then verifying the data is replaced with masked values that are realistic but fake.

**Acceptance Scenarios**:

1. **Given** a PostgreSQL database with a table `users` containing `email` and `full_name`, **When** I run `dbmask apply --config config.hcl`, **Then** the `email` and `full_name` columns are updated with fake data.
2. **Given** a valid configuration, **When** I run `dbmask apply` without `--force` on a database not in the allowlist, **Then** the tool refuses to run for safety.

---

### User Story 2 - Deterministic Masking for Consistent Local Development (Priority: P2)

As a Developer, I want the masking to be deterministic based on the primary key, so that when I re-run the masking tool on the same data, I get the same fake values every time.

**Why this priority**: Determinism ensures that linked records (if any) or simply the developer's mental model of the data remains consistent across masking runs.

**Independent Test**: Run the masking tool twice on the same dataset and verify that the resulting masked values are identical.

**Acceptance Scenarios**:

1. **Given** a row with PK `123`, **When** I mask it twice, **Then** the value generated for `first_name` is exactly the same in both runs.
2. **Given** multiple workers are used for parallelism, **When** masking is performed, **Then** the output remains identical to a single-worker run.

---

### User Story 3 - Predictable Testing with Test Fakers (Priority: P3)

As a QA Engineer, I want to use special "test fakers" that return predictable values, so that I can write exact assertions in E2E tests without dealing with random data.

**Why this priority**: Essential for robust, non-brittle automated testing of the masking logic itself.

**Independent Test**: Configure a column to use a `test_name` provider and verify that the database contains the exact expected string (e.g., "Test User 1").

**Acceptance Scenarios**:

1. **Given** a config using `faker { provider = "test_name" }`, **When** I run `dbmask apply`, **Then** the column is populated with predictable test values.

---

### Edge Cases

- **Missing Primary Key**: The tool must fail if a table being masked does not have a primary key, as it's required for deterministic seeding.
- **Unknown HCL Blocks**: The tool must fail hard if the configuration file contains unknown blocks or fields to prevent accidental misconfiguration.
- **Database Connection Failure**: Graceful exit with a clear error message if the PostgreSQL database is unreachable.
- **Unused Configuration**: If a table or column is defined in HCL but doesn't exist in the database, the tool should fail (or at least warn/fail based on "fail hard" rule).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST parse declarative HCL configuration files for masking rules.
- **FR-002**: System MUST support `faker`, `template`, `constant`, and `null` generator types.
- **FR-003**: System MUST provide a compile-time plugin-style registry for generators and faker providers.
- **FR-004**: System MUST ensure deterministic output by seeding generators with `MD5(table_name + ":" + pk_value)`.
- **FR-005**: System MUST require a Primary Key for every table being masked and use `ORDER BY pk` for processing.
- **FR-006**: System MUST only support PostgreSQL databases using the `pgx` driver.
- **FR-007**: System MUST implement a safety check that refuses execution unless the database is in an allowlist or `--force` is provided.
- **FR-008**: System MUST support parallel execution without affecting the determinism of the output.
- **FR-009**: System MUST include dedicated test-only faker providers for deterministic E2E testing.
- **FR-010**: System MUST use prepared statements for all database updates.

### Key Entities

- **Configuration**: The HCL file defining which tables and columns to mask and with which generators.
- **Generator**: An implementation of the `Generate(ctx RowContext)` interface that produces masked data.
- **Faker Provider**: A specific category of fake data (e.g., name, email) within the faker generator.
- **Registry**: The central, explicit registration point for all available generators and providers.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Tool successfully masks a 10,000-row table in under 30 seconds (excluding network latency).
- **SC-002**: 100% of masked values are deterministic across multiple runs on the same hardware and across different worker counts.
- **SC-003**: All E2E tests pass using `testcontainers-go` with exact value assertions.
- **SC-004**: Codebase has 0 linting errors and 0 build failures.
- **SC-005**: Registry allows adding a new faker provider with only Go code changes and registration, without touching core masking logic.