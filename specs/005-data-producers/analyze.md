# Analysis: Data Producers & Row Sources

## Architectural Refinements

1. **Producer Interface**:
   - `internal/db.DB` is for generic DB access.
   - `internal/producer.Producer` will wrap a `db.DB` and provide a stream of PKs.
   - Interface:
     ```go
     type Producer interface {
         EstimateCount(ctx context.Context) (int64, error)
         FetchRows(ctx context.Context, columns []string, limit int) (db.Rows, error)
         Metadata() ProducerMetadata
     }
     ```

2. **Registry Pattern**:
   - Similar to generators, producers will have a compile-time registry.
   - `NewProducer(ctx context.Context, db db.DB, cfg config.Source) (Producer, error)`.

3. **Query Safety Enforcement**:
   - The `query` producer must explicitly check the SQL for `ORDER BY` and ensure it's a `SELECT`.
   - We'll use simple string parsing or a regex for now to avoid a full SQL parser dependency, as per constraints.

4. **Engine Decoupling**:
   - The engine currently calls `e.DB.FetchRows` with table name and PK.
   - It will now call `producer.FetchRows`. The producer knows the table name and ordering logic.

5. **Dry-Run & Plan**:
   - `plan` will call `producer.EstimateCount`.
   - `dry-run` will behave exactly like apply, just with the engine's `DryRun` flag set.

## Applied Suggestions (Automatic)

- **Default Source**: If no `source` block is present, the engine will default to the `table` producer.
- **PK Consistency**: The engine still holds the `pk` name from the `table` block and passes it to the producer if needed.

## Risks & Mitigations

- **Risk**: Custom queries returning duplicate PKs.
- **Mitigation**: Warn in docs; potentially add a runtime check if it's cheap (tracking first N rows).
- **Risk**: Complex queries failing in `EstimateCount`.
- **Mitigation**: Producers should provide a safe fallback or a generic count if estimation fails.
