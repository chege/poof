# Feature Specification: Opinionated, Comfortable, Autonomous CLI

**Feature Branch**: `003-comfortable-cli`  
**Created**: 2026-02-07  
**Status**: Draft  
**Input**: User description: "SRS: poof — Opinionated, Comfortable, Autonomous CLI VERSION: 2.2"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Verify Environment Readiness (Priority: P1)

As a Database Administrator, I want to run a "doctor" command that checks my configuration, database connectivity, and safety settings, so that I have full confidence before executing any data masking operations.

**Why this priority**: Essential for the "comfort" and "trust" pillars of the project. It provides a non-destructive way to verify readiness.

**Independent Test**: Can be fully tested by running `poof doctor` and verifying it correctly identifies valid/invalid configs, reachable/unreachable databases, and allowlist status.

**Acceptance Scenarios**:

1. **Given** a valid config and reachable database, **When** I run `poof doctor`, **Then** I see a series of ✓ symbols and a final "PASS" message.
2. **Given** a database not in the allowlist, **When** I run `poof doctor`, **Then** I see an error ✗ for the safety check and a final "FAIL" message.

---

### User Story 2 - Rapid Configuration Setup (Priority: P2)

As a New User, I want to initialize a standard configuration file with comments and examples, so that I can quickly start using the tool without reading extensive documentation.

**Why this priority**: Reduces friction for new users and enforces the "opinionated" defaults from the start.

**Independent Test**: Run `poof init`, verify `poof.hcl` is created, and ensure `poof validate` passes on the generated file.

**Acceptance Scenarios**:

1. **Given** an empty directory, **When** I run `poof init`, **Then** a `poof.hcl` file is created with valid syntax and commented examples.
2. **Given** I want more detail, **When** I run `poof init --explain`, **Then** the generated config includes detailed inline annotations for every block and attribute.

---

### User Story 3 - Preview Masking Impact (Priority: P3)

As a Developer, I want to see a preview of how my data will be changed before I apply the rules, so that I can catch potential misconfigurations or unwanted side effects.

**Why this priority**: Builds trust by showing exactly what will happen.

**Independent Test**: Run `poof plan` and verify it displays a small sample (N rows) of before/after values for masked columns.

**Acceptance Scenarios**:

1. **Given** a valid config, **When** I run `poof plan`, **Then** I see a summary of tables to be masked and a diff preview for the first few rows.
2. **Given** I am in a non-TTY environment, **When** I run `poof plan`, **Then** the output contains no ANSI color codes but preserves the symbol indicators.

---

### User Story 4 - Streamlined Application (Priority: P4)

As an Automated System, I want to apply masking rules without any interactive prompts, so that the tool can be used in CI/CD pipelines or scheduled tasks.

**Why this priority**: Supports the "Autonomous" and "YOLO" execution mode.

**Independent Test**: Run `poof apply --yes` and verify it completes without waiting for input, even if a plan would normally be shown.

**Acceptance Scenarios**:

1. **Given** a valid environment, **When** I run `poof apply`, **Then** it validates, prints a plan, and then applies changes without asking for confirmation.
2. **Given** the `--yes` flag, **When** I run `poof apply --yes`, **Then** it skips the plan summary and immediately proceeds to application after validation.

---

### Edge Cases

- **Non-TTY Environment**: The tool must detect when it's not a TTY and disable colors automatically.
- **Database Connection Dropout**: If connectivity is lost during `doctor` or `plan`, it must report it as a clear failure with ✗.
- **Malformed HCL in Init**: `init` should never produce malformed HCL.
- **No Primary Key**: `plan` and `apply` must fail immediately if a table lacks a primary key.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a `doctor` command that verifies config, connectivity, safety, and providers.
- **FR-002**: System MUST provide an `init` command that generates a documented `poof.hcl` file.
- **FR-003**: System MUST provide a `validate` command that checks HCL syntax and schema (unknown fields/blocks).
- **FR-004**: System MUST provide a `plan` command that shows a summary of changes and a limited diff (first N rows).
- **FR-005**: System MUST provide an `apply` command that incorporates validation and planning before execution.
- **FR-006**: System MUST use symbols (✓, !, ✗) for status reporting.
- **FR-007**: System MUST support `--no-color` and detect TTY status for color output.
- **FR-008**: System MUST NOT use interactive prompts or require elevation (sudo).
- **FR-009**: System MUST support a `--yes` flag for `apply` to skip the plan summary.

### Key Entities

- **Doctor Report**: A structured summary of environment readiness checks.
- **Masking Plan**: A non-persistent summary of intended database changes, including a small sample diff.
- **Standard Config**: The opinionated HCL template generated by `init`.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can generate a working configuration file in under 2 seconds using `poof init`.
- **SC-002**: `poof doctor` completes all checks and reports status in under 5 seconds for local databases.
- **SC-003**: Tool provides consistent output in both TTY and non-TTY environments (with/without colors).
- **SC-004**: 0% of commands require `sudo` or interactive user confirmation.
- **SC-005**: `poof plan` output is limited to a small, readable number of sample rows (e.g., max 5 per table).