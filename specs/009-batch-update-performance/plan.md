# Implementation Plan: Batch Update Performance

**Branch**: `009-batch-update-performance` | **Date**: Saturday, February 7, 2026 | **Spec**: [specs/009-batch-update-performance/spec.md](specs/009-batch-update-performance/spec.md)
**Input**: Feature specification from `specs/009-batch-update-performance/spec.md`

## Summary

Implement high-performance batch updates for the masking engine. Instead of executing one `UPDATE` per row, Poof will collect masked values into memory and apply them in bulk using PostgreSQL's `UPDATE ... FROM (VALUES ...)` syntax. On any batch failure (e.g., unique constraint violation), the system will automatically fall back to row-by-row processing for that specific batch to maintain robustness.

## Technical Context

**Language/Version**: Go 1.25.7
**Primary Dependencies**: 
- `github.com/jackc/pgx/v5`
- `internal/engine`
- `internal/config`
**Performance Goals**: 10x-50x increase in masking throughput for large tables.
**Default Batch Size**: 500 rows.

## Constitution Check

- **Safe-by-default**: Pass. Transactions and row-by-row fallback ensure no data is lost on batch failure.
- **Deterministic**: Pass. Batching does not change the generated values or the seeding logic.
- **No Magic**: Pass. Batch size is explicit and configurable.

## Project Structure

### Documentation (this feature)

```text
specs/009-batch-update-performance/
├── plan.md              # This file
├── spec.md              # Feature specification
└── checklists/
    └── requirements.md
```

### Source Code

```text
internal/
├── config/
│   └── models.go       # Update: Add BatchSize to Config
├── engine/
│   ├── engine.go       # Update: Implement batching logic in writeResults
│   └── batch.go        # New: Helper for building bulk SQL
```

## Research (Phase 0)

### 1. PostgreSQL Bulk Update Syntax
The most efficient method without temporary tables:
```sql
UPDATE target_table AS t
SET 
    col1 = v.col1,
    col2 = v.col2
FROM (VALUES 
    (1, 'val1_a', 'val2_a'),
    (2, 'val1_b', 'val2_b')
) AS v(pk, col1, col2)
WHERE t.pk = v.pk;
```

### 2. Implementation in Poof
- **Buffer**: `writeResults` will use a slice `[]rowData` as a buffer.
- **Trigger**: When `len(buffer) == BatchSize` or `outputChan` is closed.
- **SQL Building**: A helper will generate the `UPDATE ... FROM (VALUES ...)` string with appropriate placeholders ($1, $2, ...).
- **Fallback**: 
    ```go
    if err := e.applyBatch(ctx, tx, tableCfg, columnNames, buffer, results); err != nil {
        slog.Warn("batch update failed, falling back to individual updates", "error", err)
        for _, row := range buffer {
            e.retryUpdate(ctx, tx, singleRowQuery, tableCfg, columnNames, generators, &row, maxRetries, results)
        }
    }
    ```

## Data Model (Phase 1)

### Configuration Extensions (`internal/config/models.go`)
```go
type Config struct {
    // ...
    BatchSize int `toml:"batch_size"`
}
```

## Detailed Tasks (Phase 2 Preview)

1.  **Config**: Add `BatchSize` to `Config` struct and default to 500 in `LoadConfig`.
2.  **Engine**: Update `Engine` struct to include `BatchSize`.
3.  **SQL**: Implement `buildBatchUpdateQuery` in `internal/engine/batch.go`.
4.  **Engine**: Refactor `writeResults` to use a buffer and call batch update.
5.  **Robustness**: Implement the fallback loop for failed batches.
6.  **Verification**: Add integration test with large dataset to verify performance and fallback.