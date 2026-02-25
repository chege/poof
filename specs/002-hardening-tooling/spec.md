# Feature Specification: poof — Autonomous Hardening & Tooling Phase

**Feature Branch**: `002-hardening-tooling`  
**Created**: 2026-02-07  
**Status**: Draft  
**Input**: User description: "SRS: poof — Autonomous Hardening & Tooling Phase VERSION: 2.1"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Standardized Development Workflow (Priority: P1)

As a developer, I want to use a set of opinionated tasks (via Taskfile.dev) to build, test, and verify the tool, so that the development lifecycle is consistent, automated, and guarded against common mistakes.

**Why this priority**: This is a foundational requirement for autonomous development and ensures quality gates are strictly enforced.

**Independent Test**: Can be fully tested by running `task check` and verifying it executes `fmt`, `vet`, and `test` in the correct sequence, failing if any step fails.

**Acceptance Scenarios**:

1. **Given** a clean repository, **When** I run `task build`, **Then** the `poof` binary is produced without errors.
2. **Given** a codebase with linting or test failures, **When** I run `task check`, **Then** the process stops and reports the exact failure.

---

### User Story 2 - Expanded Data Masking Capabilities (Priority: P2)

As a user, I want to have a broader set of deterministic faker providers (username, company_name, phone_number, ipv4_address, short_text), so that I can mask a wider variety of sensitive data types in my database.

**Why this priority**: Enhances the utility of the tool by supporting more common database fields while maintaining the core principle of determinism.

**Independent Test**: Can be tested by configuring a table with the new providers in `poof.toml` and verifying that the database contains realistic but fake data that remains identical across multiple runs.

**Acceptance Scenarios**:

1. **Given** a column configured with `gen "faker" { provider = "ipv4_address" }`, **When** I apply the masking, **Then** the column is populated with valid IPv4 addresses.
2. **Given** multiple runs of the same configuration, **When** checking the `username` provider output for the same row, **Then** the value is identical every time.

---

### User Story 3 - Robust and Hardened Implementation (Priority: P3)

As a system, I want to have explicit error handling, strict validation of configuration, and clear package boundaries, so that the tool is safer to use, easier to maintain, and provides clear feedback on failure.

**Why this priority**: Improves the "hardened" aspect of the tool, reducing the risk of silent failures or misconfiguration.

**Independent Test**: Can be tested by providing an TOML configuration with unknown fields or blocks and verifying that the tool fails hard with an explicit error message.

**Acceptance Scenarios**:

1. **Given** an TOML config with an unknown attribute, **When** the tool loads it, **Then** it fails immediately with a descriptive error.
2. **Given** a table without a primary key, **When** masking is attempted, **Then** the tool provides a clear error explaining the PK requirement.

---

### Edge Cases

- **Taskfile missing**: The system should rely on `task` being available as a prerequisite.
- **Malformed TOML**: System must handle TOML syntax errors gracefully but firmly (fail hard).
- **Parallelism with new providers**: Ensure `ipv4_address` and `short_text` are correctly seeded and thread-safe.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST implement a `Taskfile.yml` containing at least: `build`, `test`, `fmt`, `vet`, `check`, `init`, `plan`, `apply`.
- **FR-002**: System MUST add five new deterministic faker providers: `username`, `company_name`, `phone_number` (locale-agnostic), `ipv4_address`, `short_text` (fixed-length).
- **FR-003**: System MUST provide test fakers for all new providers that return predictable values for exact assertions.
- **FR-004**: System MUST include a deterministic demo dataset for examples and documentation.
- **FR-005**: System MUST harden the codebase by removing dead code and clarifying package boundaries.
- **FR-006**: System MUST improve error messages to be explicit and actionable.

### Key Entities

- **Task**: An opinionated, named command in `Taskfile.yml` that wraps development operations.
- **Hardened Engine**: The updated masking logic with improved validation and explicit error paths.
- **Demo Data**: A predefined set of SQL records used to demonstrate and test masking.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All tasks in `Taskfile.yml` execute successfully in the target environment.
- **SC-002**: 100% of new faker providers are verified to be deterministic across parallel workers.
- **SC-003**: E2E tests pass for the full masking flow using `testcontainers-go`, including verification of new providers.
- **SC-004**: Tool fails with a non-zero exit code and clear message when encountering unknown TOML configuration.
- **SC-005**: Codebase is free of `init()` side effects beyond explicit registration, as per core principles.