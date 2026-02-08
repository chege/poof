# Feature Specification: Batch Update Performance

**Feature Branch**: `009-batch-update-performance`  
**Created**: Saturday, February 7, 2026  
**Status**: Draft  
**Input**: User description: "Implement update batching. Instead of updating rows one-by-one, collect masked values into groups (e.g., 500 rows) and use UPDATE ... FROM (VALUES ...) or temporary tables to apply changes in bulk. This will provide a 10x-50x performance boost for large datasets."

## User Scenarios & Testing

### User Story 1 - High-Performance Masking for Large Tables (Priority: P1)

As a data engineer, I want the masking process to be significantly faster when dealing with tables containing millions of rows, so that I can complete masking jobs within short maintenance windows.

**Why this priority**: Scalability is critical for enterprise use cases where single-row updates are prohibitively slow.

**Independent Test**: Can be tested by masking a table with 100,000+ rows and comparing the execution time of batch mode vs. the previous single-row mode.

**Acceptance Scenarios**:

1. **Given** a table with 100,000 rows, **When** `poof apply` is executed with a batch size of 500, **Then** the total masking time is at least 10x faster than single-row updates.
2. **Given** a default configuration, **When** masking starts, **Then** the system automatically uses a reasonable default batch size (e.g., 500) without user intervention.

---

### User Story 2 - Graceful Degradation on Unique Violations (Priority: P2)

As an operator, I want the system to handle `UNIQUE` constraint violations within a batch by falling back to row-by-row updates for that specific batch, so that I don't lose the performance of batching for the rest of the table.

**Why this priority**: Ensures robustness. A single collision shouldn't invalidate the performance gains for an entire dataset.

**Independent Test**: Can be tested by intentionally introducing a unique constraint collision in one row of a batch and verifying that the other rows in that batch are still masked and the job continues.

**Acceptance Scenarios**:

1. **Given** a batch of 500 rows where 1 row violates a `UNIQUE` constraint, **When** the batch update fails, **Then** the system automatically retries those 500 rows individually to identify and handle the specific failure.

---

### Edge Cases

- **Small Tables**: How does the system handle tables smaller than the batch size?
- **Memory Pressure**: What is the impact of large batch sizes on the CLI's memory consumption?
- **Transaction Boundaries**: Does a batch failure affect the atomicity of the entire table masking operation?

## Requirements

### Functional Requirements

- **FR-001**: System MUST collect masked values into memory buffers before sending updates to the database.
- **FR-002**: System MUST use a high-performance bulk update method (e.g., `UPDATE ... FROM (VALUES ...)` or temporary tables).
- **FR-003**: System MUST support a configurable batch size via the configuration file or environment variables.
- **FR-004**: System MUST default to a batch size of 500 rows if not specified.
- **FR-005**: System MUST detect batch-level failures (e.g., constraint violations) and automatically fall back to single-row processing for that specific batch.
- **FR-006**: System MUST maintain determinism and transactional integrity during batch processing.

### Assumptions

- **AS-001**: `UPDATE ... FROM (VALUES ...)` is supported by the target PostgreSQL versions (14+).
- **AS-002**: Performance gains will vary based on network latency and database hardware.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Masking throughput increases by at least 10x for tables with >50,000 rows compared to single-row updates.
- **SC-002**: CLI memory usage remains below 200MB even when processing batches of 1,000 rows with large text columns.
- **SC-003**: 100% of rows are correctly masked even when a batch requires fallback due to a constraint violation.
- **SC-004**: Zero data loss or corruption occurs across entire tables when using batching.