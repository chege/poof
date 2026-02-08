# Feature Specification: Config Pre-flight Validation

**Feature Branch**: `011-config-preflight-validation`  
**Created**: Sunday, February 8, 2026  
**Status**: Draft  
**Input**: User description: "🛡 Config Validation CLI (Pre-flight Check) Right now, you find out your config is invalid when you try to run apply or plan. * Recommendation: Enhance poof validate to perform deep semantic validation, not just syntax checking. * Checks: Verify database connectivity, check if all tables/columns exist in the DB, validate generator parameters (e.g., regex compilation), and ensure salt is present for production environments."

## User Scenarios & Testing

### User Story 1 - Semantic Configuration Validation (Priority: P1)

As a developer, I want to verify that my masking configuration is semantically correct before running it against a database, so that I can avoid runtime failures caused by typos or missing tables.

**Why this priority**: Fundamental safety and DX requirement. Reduces frustration by catching errors early.

**Independent Test**: Can be tested by running `poof validate` with a configuration referencing a non-existent table and verifying that it fails with a specific error message.

**Acceptance Scenarios**:

1. **Given** a `poof.toml` that references a table not present in the database, **When** `poof validate` is executed, **Then** it returns a non-zero exit code and explicitly lists the missing table.
2. **Given** a valid `poof.toml` and an active database connection, **When** `poof validate` is executed, **Then** it returns exit code 0 and reports "Configuration is valid and schema verified."

---

### User Story 2 - Generator Parameter Verification (Priority: P1)

As an operator, I want the validation command to check that all generator parameters are valid (e.g., valid Go templates, valid faker providers), so that I don't discover these errors mid-way through a long-running masking job.

**Why this priority**: Prevents partial job failures that are expensive to recover from.

**Independent Test**: Can be tested by providing an invalid Go template string in a `template` generator and running `poof validate`.

**Acceptance Scenarios**:

1. **Given** a column using a `template` generator with a syntax error, **When** `poof validate` is run, **Then** the validation fails with a "template parse error" for that specific column.
2. **Given** an unknown `faker` provider name, **When** `poof validate` is run, **Then** the validation fails listing the unknown provider.

---

### User Story 3 - Production Readiness Check (Priority: P2)

As a security engineer, I want the validation tool to enforce safety standards for production environments, such as requiring a global salt, so that our masking remains secure and deterministic across runs.

**Why this priority**: Hardens the security posture of the tool for enterprise use.

**Independent Test**: Can be tested by running validation against a production database without a `salt` defined in the `[safety]` block.

**Acceptance Scenarios**:

1. **Given** a database environment that is likely production (based on allowlist or name), **When** `poof validate` is run without a global `salt`, **Then** it issues a warning or error (configurable strictness).

---

### Edge Cases

- **Disconnected DB**: How does `poof validate` behave if it cannot connect to the database? (It should report connectivity failure as a distinct validation error).
- **Partial Permissions**: What if the tool can connect but doesn't have permission to read `information_schema`?
- **Ambiguous Environments**: How does validation handle multiple environments in a single config? (It should validate the selected or default one).

## Requirements

### Functional Requirements

- **FR-001**: System MUST perform a database connectivity check during validation.
- **FR-002**: System MUST verify that all tables defined in the config exist in the target database.
- **FR-003**: System MUST verify that all columns defined in the config exist in their respective tables.
- **FR-004**: System MUST validate the syntax of all `template` generators.
- **FR-005**: System MUST verify that all `faker` providers used are registered and available.
- **FR-006**: System MUST check for the presence of a non-empty `salt` if the database name matches specific patterns or is marked as high-risk.
- **FR-007**: System MUST provide clear, actionable error messages including table and column names where applicable.
- **FR-008**: System MUST support a `--strict` mode that treats all warnings as errors.

### Assumptions

- **AS-001**: Validation requires an active network connection to the target database.
- **AS-002**: The tool will use the existing `GetTableColumns` and `GetAllTables` infrastructure.

## Success Criteria

### Measurable Outcomes

- **SC-001**: 100% of runtime errors related to schema mismatches or invalid generator parameters are caught by `poof validate`.
- **SC-002**: Configuration validation for a schema with 100 tables completes in under 2 seconds.
- **SC-003**: All validation errors return the correct structured exit code (`2` for Config Error).
- **SC-004**: User can identify exactly which line or field in the TOML caused the semantic failure from the output.