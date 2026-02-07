# Analysis: poof-cli

## Architectural Refinements

1. **Deterministic Seeding**:
   - To ensure absolute determinism in a parallel environment, the `RowContext` will carry the calculated seed.
   - Each generator `Generate` call will be passed this `RowContext`.
   - For `faker`, we will use a local random source seeded with the calculated seed for *every* row, rather than a shared global source. This prevents race conditions and ensures output is independent of worker count or execution order.

2. **HCL Validation**:
   - We will use `hcl.Diagnostics` to report all errors.
   - "Fail hard" on unknown fields will be enforced using `hcl.Body.Content` and checking for unhandled attributes.

3. **Database Efficiency**:
   - While using one transaction per table, we will use a `PREPARE` statement for the update query and execute it in a loop (or batch) for all rows.
   - The `SELECT` query must include `ORDER BY pk` to ensure consistent processing order, although seeding makes this less critical for the values themselves, it's good for progress tracking and predictability.

4. **Safety Verification**:
   - The allowlist check will be performed against the `database` name retrieved from the `pgx` connection config.
   - If `apply` is called, the tool will first query the database name to verify it matches the allowlist.

5. **Test Faker Strategy**:
   - Test fakers will be registered under a special namespace or simply as providers that are always available but documented for testing use.
   - Example: `faker { provider = "test_name" }`.

## Potential Risks & Mitigations

- **Risk**: Memory usage for large tables in a single transaction.
- **Mitigation**: We will stream rows using `pgx.Rows` and update them row-by-row or in small batches. We won't load the entire table into memory.
- **Risk**: Concurrent updates causing deadlocks.
- **Mitigation**: Since we are masking the whole table and using `ORDER BY pk` for selection, but updates are independent, we should be careful. However, with one transaction per table, we shouldn't have concurrent transactions touching the same table from *this* tool.

## Applied Suggestions (Automatic)

- **Registry**: I will implement a `Register(name string, gen Generator)` pattern.
- **Context**: `RowContext` will include `TableName`, `ColumnName`, `PrimaryKeyValue`, and `Seed`.
- **Boilerplate**: I will immediately fix the empty `main.go` and `root.go`.
