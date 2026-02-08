# Feature Specification: Poof Core Enhancements

**Feature Branch**: `008-poof-core-v1-enhancements`  
**Created**: Saturday, February 7, 2026  
**Status**: Draft  
**Input**: User description: "The high and medium priority ones (Structured Exit Codes, Multiple DB Blocks, Locale Support, Generator Composition)"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - CI/CD Integration with Structured Exit Codes (Priority: P1)

As a DevOps engineer, I want the Poof CLI to return specific exit codes for different failure modes, so that my CI/CD pipelines can automatically decide whether to retry the job (on network errors) or alert the team (on validation errors).

**Why this priority**: Critical for automation and reliable integration into production pipelines.

**Independent Test**: Can be tested by intentionally triggering different error types (invalid config, connection failure, unique violation) and verifying the process exit code using `echo $?`.

**Acceptance Scenarios**:

1. **Given** an invalid TOML configuration, **When** `poof apply` is run, **Then** the exit code is `2` (Configuration Error).
2. **Given** an unreachable database, **When** `poof apply` is run, **Then** the exit code is `3` (Connection Error).
3. **Given** a masking job that completes with row failures, **When** it finishes, **Then** the exit code is `4` (Masking Failures occurred).

---

### User Story 2 - Multi-Environment Configuration (Priority: P1)

As a developer, I want to define my `local`, `staging`, and `prod` database connection details in a single `poof.toml` file, so that I don't have to manage multiple configuration files for different environments.

**Why this priority**: Major quality-of-life improvement for users managing complex deployments.

**Independent Test**: Can be tested by creating a config with multiple named database blocks and running `poof apply --env staging`.

**Acceptance Scenarios**:

1. **Given** a config with `[databases.local]` and `[databases.staging]`, **When** `poof apply --env staging` is run, **Then** it connects to the staging DSN.
2. **Given** no `--env` flag is provided, **When** `poof apply` is run, **Then** it defaults to the first defined environment or a block marked as `default`.

---

### User Story 3 - Localized Masking Data (Priority: P2)

As a QA engineer for a German application, I want my masked user data to look like real German names and addresses, so that my localized search and validation features can be tested effectively.

**Why this priority**: Necessary for international applications where US-centric data is insufficient for testing.

**Independent Test**: Can be tested by setting `locale = "de_DE"` in the configuration and verifying that `faker` generates non-English names.

**Acceptance Scenarios**:

1. **Given** a global locale setting of `de_DE`, **When** the `faker` generator is used for names, **Then** the generated names follow German conventions.
2. **Given** a specific column has a locale override, **When** generated, **Then** it uses the specific locale instead of the global default.

---

### User Story 4 - Combined Generator Output (Priority: P2)

As a data engineer, I want to combine static text with dynamic counters or faker data (e.g., "USER-" + counter), so that I can create complex, unique identifiers that meet my application's requirements.

**Why this priority**: Increases the flexibility of generators without requiring custom Go code.

**Independent Test**: Can be tested by configuring a column to use a "composed" generator and verifying the output format.

**Acceptance Scenarios**:

1. **Given** a column configured with a template like `"MKT-{{ .counter }}"`, **When** masked, **Then** the output is a string starting with "MKT-" followed by the deterministic counter value.

---

### Edge Cases

- **Environment Collisions**: What happens if the same environment name is defined twice?
- **Invalid Locales**: How does the system handle a requested locale that the faker library doesn't support?
- **Circular Composition**: Does the system prevent templates from referencing themselves in a way that causes an infinite loop?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST define and return a set of distinct exit codes (0: Success, 1: General, 2: Config, 3: Connection, 4: Partial Failure).
- **FR-002**: System MUST support a `[databases]` table in TOML allowing multiple named environments.
- **FR-003**: System MUST provide an `--env` flag to select the active database environment.
- **FR-004**: System MUST support a global and per-column `locale` setting.
- **FR-005**: System MUST integrate a localization-aware faker library (supporting at least `en`, `de`, `fr`, `es`).
- **FR-006**: System MUST allow generators to be composed using Go-template syntax within the configuration.
- **FR-007**: System MUST remain backward compatible with single `[database]` block configurations (deprecated but functional).

### Assumptions

- **AS-001**: The existing `faker` providers will be mapped to the new localized implementation.
- **AS-002**: Template-based composition will use the `template` generator as the underlying mechanism.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of CLI errors return a non-zero exit code corresponding to their category.
- **SC-002**: Users can switch between `staging` and `local` masking targets using only a command-line flag.
- **SC-003**: Generated data matches the phonetic and orthographic conventions of at least 4 supported locales.
- **SC-004**: Complex ID strings (e.g., "REF-1001") can be generated without any custom Go code.