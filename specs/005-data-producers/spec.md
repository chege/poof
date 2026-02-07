# Feature Specification: Data Producers & Row Sources

**Feature Branch**: `005-data-producers`  
**Created**: 2026-02-07  
**Status**: Draft  
**Input**: User description: "SRS: poof — Data Producers & Row Sources VERSION: 5.0"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Standard Table Masking (Priority: P1)

As a user, I want to mask an entire table using the default source configuration, so that I can easily protect all records in a standard data layout.

**Why this priority**: This is the existing behavior and remains the foundation of the tool.

**Independent Test**: Can be tested by omitting the `source` block in a `table` configuration and verifying that all rows are selected and masked.

**Acceptance Scenarios**:

1. **Given** a table `users` with PII, **When** I run `poof apply` with no `source` block, **Then** all rows in `users` are selected via `SELECT pk FROM users ORDER BY pk` and masked.

---

### User Story 2 - Masking Data Subsets via Views (Priority: P2)

As a user, I want to define a database view to filter the data I want to mask, while having the updates applied to the underlying table, so that I can target specific data segments or tenants.

**Why this priority**: Supports real-world scenarios where only a portion of the data should be masked (e.g., specific countries, active users).

**Independent Test**: Create a view `active_users`, configure `table "users" { source { type = "view" name = "active_users" } }`, and verify only rows in the view are masked.

**Acceptance Scenarios**:

1. **Given** a view `v_active_users` on table `users`, **When** I configure it as a source, **Then** only rows returned by the view are processed, but the `UPDATE` statements target the `users` table.

---

### User Story 3 - Targeted Masking via Custom Query (Priority: P3)

As an advanced user, I want to provide a custom SQL query to select exactly which rows should be masked, so that I can perform incremental or highly specific masking operations.

**Why this priority**: Provides maximum flexibility for complex row selection requirements.

**Independent Test**: Configure a `table` with `source { type = "query" sql = "..." }` and verify the selected rows are masked according to the config.

**Acceptance Scenarios**:

1. **Given** a custom SELECT query that returns PKs, **When** I run `poof apply`, **Then** the engine uses those PKs to perform the masking.
2. **Given** a query without an `ORDER BY pk` clause, **When** validation runs, **Then** the tool refuses to proceed to ensure determinism.

---

### Edge Cases

- **Non-updatable Views**: System must reject views that cannot be used for updates (though we update the base table, the view must still be valid).
- **Query Cardinality**: Custom queries that join tables must be checked to ensure they don't return duplicate PKs (which would cause multiple updates per row).
- **Missing PK in Source**: Both `view` and `query` sources must return the primary key column specified in the `table` block.
- **Unknown Source Type**: Fails validation immediately.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Introduce a `Producer` interface that abstracts row selection.
- **FR-002**: Support `table` (default), `view`, and `query` producer types.
- **FR-003**: The `table` block in HCL MUST support an optional `source` block.
- **FR-004**: Producers MUST yield row identifiers (PKs) in a deterministic order.
- **FR-005**: The `view` producer MUST map rows to a base table for updates.
- **FR-006**: The `query` producer MUST be a `SELECT` statement returning the PK and MUST include `ORDER BY pk`.
- **FR-007**: Validation MUST fail for queries that modify cardinality or hide ordering.
- **FR-008**: Producers MUST participate in `plan`, `dry-run`, and `doctor` commands.
- **FR-009**: Producers MUST be registered at compile-time via an explicit registry.

### Key Entities

- **Producer**: The component responsible for identifying rows to be masked.
- **Source Configuration**: The HCL block defining the producer type and its parameters.
- **Engine**: The core logic that consumes rows from a producer and applies generators.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: System successfully masks rows selected via all three producer types.
- **SC-002**: Safety checks (allowlist, dry-run) remain fully functional for all producer types.
- **SC-003**: Custom queries without `ORDER BY` are rejected with a clear error message.
- **SC-004**: Adding a new producer in code requires zero changes to the engine's core masking logic.
- **SC-005**: `poof plan` correctly estimates row counts for views and custom queries.