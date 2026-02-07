# Analysis: Release Readiness, Safety & Extensible Database Design

## Architectural Refinements

1. **Minimal Database Abstraction**:
   - The core engine should not know about SQL dialects or driver-specific types.
   - `internal/db.DB` interface will act as the bridge.
   - `internal/db.Rows` and `internal/db.Tx` interfaces will abstract row iteration and transactions.

2. **Implicit Backend Selection**:
   - We will implement a `Connect(dsn string)` function in `internal/db`.
   - This function will inspect the connection string scheme (e.g., `postgres://`) and route to the appropriate factory.
   - Unsupported schemes will return a clear error: `"unsupported database scheme: mysql"`.

3. **Strict Dry-Run enforcement**:
   - In `DryRun` mode, the `masker.Engine` will perform all data generation steps but will skip the `tx.Commit()` or use a read-only transaction if supported by the backend.
   - To be absolutely safe, the PostgreSQL implementation can use `Rollback()` at the end of a dry-run session even if no changes were made.

4. **Efficient Estimates**:
   - For `dbmask plan`, we need row counts. 
   - `SELECT count(*)` is slow on large tables.
   - We will implement `EstimateRowCount` in the interface. For PostgreSQL, we'll try to use `pg_class.reltuples` for a fast, near-instant estimate.

## Applied Suggestions (Automatic)

- **Package Reorganization**: I will move the existing `postgresql.go` into `internal/db/postgres/client.go` to enforce package boundaries.
- **Interface Implementation**: The new interfaces will be designed to be minimal to avoid the "Fat Interface" anti-pattern.

## Risks & Mitigations

- **Risk**: Information leakage between DB backends.
- **Mitigation**: Use strictly generic types in the interface. Do not pass `pgx.Rows` directly; wrap it in an internal `Rows` interface.
- **Risk**: Performance hit from interface overhead.
- **Mitigation**: Minimal. Go's interface overhead is negligible compared to database I/O.
