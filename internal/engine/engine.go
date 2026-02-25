// Package engine coordinates the data masking process.
package engine

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/jackc/pgconn"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/sync/errgroup"

	"github.com/chege/poof/internal/config"
	"github.com/chege/poof/internal/db"
	"github.com/chege/poof/internal/generator"
	"github.com/chege/poof/internal/producer"
)

const maxRetries = 10

// Engine orchestrates the data masking process across multiple tables and workers.
type Engine struct {
	DB         db.DB
	Config     *config.Config
	seenValues map[string]map[any]bool
	Workers    int
	BatchSize  int
	seenMu     sync.RWMutex
	DryRun     bool
	PlanOnly   bool
}

// TableReport provides a summary of masking results for a specific table.
type TableReport struct {
	Name      string
	Columns   []string
	Types     []string // Corresponding SQL types for columns (index 0 is PK)
	Estimates int64
	Updated   int64
	Retried   int64
	Failed    int64
}

// MaskingReport provides a comprehensive summary of masking operations.
type MaskingReport struct {
	StartTime time.Time
	EndTime   time.Time
	Tables    []TableReport
	Diffs     []Diff
	Duration  time.Duration
}

// Diff represents a single row value change.
type Diff struct {
	PKValue    any
	OldValue   any
	NewValue   any
	TableName  string
	ColumnName string
}

// NewEngine creates a new masking engine with the given database client and configuration.
func NewEngine(client db.DB, cfg *config.Config, workers int) *Engine {
	if workers <= 0 {
		workers = 1
	}
	batchSize := 500
	if cfg != nil && cfg.BatchSize > 0 {
		batchSize = cfg.BatchSize
	}
	return &Engine{
		DB:         client,
		Config:     cfg,
		Workers:    workers,
		BatchSize:  batchSize,
		seenValues: make(map[string]map[any]bool),
	}
}

type rowData struct {
	pkValue   any
	oldValues []any
	newValues []any
}

// Apply executes the masking rules as defined in the configuration.
func (e *Engine) Apply(ctx context.Context) (*MaskingReport, error) {
	producer.RegisterAll()
	report := &MaskingReport{
		StartTime: time.Now(),
	}

	if !e.DryRun && !e.PlanOnly {
		if err := e.DB.SetJobState(ctx, "STARTED"); err != nil {
			slog.Warn("failed to set job state to STARTED", "error", err)
		}
	}

	for _, tableCfg := range e.Config.Tables {
		p, err := producer.NewProducer(ctx, e.DB, tableCfg.Name, tableCfg.PK, tableCfg.Source)
		if err != nil {
			return nil, fmt.Errorf("failed to create producer for %s: %w", tableCfg.Name, err)
		}

		count, err := p.EstimateCount(ctx)
		if err != nil {
			count = 0
		}

		cols := make([]string, len(tableCfg.Columns))
		for i, c := range tableCfg.Columns {
			cols[i] = c.Name
		}

		tr := TableReport{
			Name:      tableCfg.Name,
			Estimates: count,
			Columns:   cols,
		}

		if e.PlanOnly {
			report.Tables = append(report.Tables, tr)
			continue
		}

		diffs, tableResults, err := e.maskTable(ctx, tableCfg, p)
		if err != nil {
			return nil, err
		}
		tr.Updated = tableResults.Updated
		tr.Retried = tableResults.Retried
		tr.Failed = tableResults.Failed

		report.Tables = append(report.Tables, tr)
		report.Diffs = append(report.Diffs, diffs...)
	}
	report.EndTime = time.Now()
	report.Duration = report.EndTime.Sub(report.StartTime)

	if !e.DryRun && !e.PlanOnly {
		if err := e.DB.SetJobState(ctx, "COMPLETED"); err != nil {
			slog.Warn("failed to set job state to COMPLETED", "error", err)
		}
	}

	return report, nil
}

func (e *Engine) maskTable(ctx context.Context, tableCfg config.Table, p producer.Producer) ([]Diff, TableReport, error) {
	results := TableReport{Name: tableCfg.Name}
	columnNames, generators, err := e.setupGenerators(tableCfg)
	if err != nil {
		return nil, results, err
	}

	// Fetch types for bulk update casting
	tableCols, err := e.DB.GetTableColumns(ctx, tableCfg.Name)
	if err == nil {
		typeMap := make(map[string]string)
		for _, tc := range tableCols {
			typeMap[tc.Name] = tc.DataType
		}
		// Index 0 is PK
		results.Types = append(results.Types, typeMap[tableCfg.PK])
		for _, cn := range columnNames {
			results.Types = append(results.Types, typeMap[cn])
		}
	}

	limit := 0
	if e.DryRun {
		limit = 5
	}

	filter := ""
	if tableCfg.Source != nil {
		filter = tableCfg.Source.Filter
	}

	rows, err := p.FetchRows(ctx, columnNames, filter, limit)
	if err != nil {
		return nil, results, fmt.Errorf("failed to fetch rows: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var tx db.Tx
	var singleUpdateQuery string
	if !e.DryRun {
		tx, err = e.DB.Begin(ctx)
		if err != nil {
			return nil, results, fmt.Errorf("failed to start transaction: %w", err)
		}
		defer func() {
			if tx != nil {
				_ = tx.Rollback(ctx)
			}
		}()

		singleUpdateQuery = e.buildUpdateQuery(tableCfg.Name, tableCfg.PK, columnNames)
	}

	g, gCtx := errgroup.WithContext(ctx)
	inputChan := make(chan rowData, e.Workers*2)
	outputChan := make(chan rowData, e.Workers*2)

	// Reader
	g.Go(func() error {
		defer close(inputChan)
		for rows.Next() {
			v, rowErr := rows.Values()
			if rowErr != nil {
				return fmt.Errorf("failed to read values: %w", rowErr)
			}
			select {
			case inputChan <- rowData{pkValue: v[0], oldValues: v[1:]}:
			case <-gCtx.Done():
				return gCtx.Err()
			}
		}
		return rows.Err()
	})

	// Workers
	for i := 0; i < e.Workers; i++ {
		g.Go(func() error {
			for {
				select {
				case <-gCtx.Done():
					return gCtx.Err()
				case row, ok := <-inputChan:
					if !ok {
						return nil
					}
					row.newValues = make([]any, len(columnNames))
					for j, colName := range columnNames {
						seedBy := tableCfg.Columns[j].SeedBy
						forceUnique := tableCfg.Columns[j].ForceUnique

						attempt := 0
						var val any
						var genErr error

						for {
							val, genErr = generators[colName].Generate(generator.NewRowContext(
								tableCfg.Name,
								colName,
								e.Config.Locale,
								e.Config.Safety.Salt,
								row.pkValue,
								row.oldValues[j],
								seedBy,
								attempt,
							))
							if genErr != nil {
								return fmt.Errorf("gen error: %w", genErr)
							}

							if !forceUnique || !e.isDuplicate(tableCfg.Name, colName, val) {
								break
							}

							attempt++
							if attempt >= maxRetries {
								slog.Warn("failed to satisfy logical uniqueness after max retries",
									"table", tableCfg.Name, "column", colName, "pk", row.pkValue)
								break
							}
							// In a real retry, we should modify the seed.
							// Currently NewRowContext creates a seed based on table+pk or value.
							// For logical uniqueness retries to work, we need a way to vary the seed.
							// I'll add an 'Attempt' field to RowContext or vary the Salt.
						}
						row.newValues[j] = val
					}
					select {
					case outputChan <- row:
					case <-gCtx.Done():
						return gCtx.Err()
					}
				}
			}
		})
	}

	// Coordinator
	go func() {
		_ = g.Wait()
		close(outputChan)
	}()

	diffs, err := e.writeResults(ctx, tableCfg, columnNames, generators, outputChan, tx, singleUpdateQuery, &results, p)
	if err != nil {
		return nil, results, err
	}

	if err := g.Wait(); err != nil {
		return nil, results, fmt.Errorf("table %s masking failed: %w", tableCfg.Name, err)
	}

	if !e.DryRun {
		if err := tx.Commit(ctx); err != nil {
			return nil, results, fmt.Errorf("commit error: %w", err)
		}
		tx = nil
	}

	return diffs, results, nil
}

func (e *Engine) isDuplicate(tableName, columnName string, value any) bool {
	key := fmt.Sprintf("%s:%s", tableName, columnName)
	e.seenMu.Lock()
	defer e.seenMu.Unlock()

	if e.seenValues[key] == nil {
		e.seenValues[key] = make(map[any]bool)
	}

	if e.seenValues[key][value] {
		return true
	}

	e.seenValues[key][value] = true
	return false
}

func (e *Engine) setupGenerators(tableCfg config.Table) ([]string, map[string]generator.Generator, error) {
	columnNames := make([]string, 0, len(tableCfg.Columns))
	generators := make(map[string]generator.Generator)

	for _, col := range tableCfg.Columns {
		columnNames = append(columnNames, col.Name)
		gen, err := generator.NewGenerator(col.Gen)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create generator for column %s: %w", col.Name, err)
		}
		generators[col.Name] = gen
	}
	return columnNames, generators, nil
}

func (e *Engine) buildUpdateQuery(tableName, pk string, columns []string) string {
	query := fmt.Sprintf("UPDATE %s SET ", tableName)
	for i, col := range columns {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", col, i+2)
	}
	query += fmt.Sprintf(" WHERE %s = $1", pk)
	return query
}

func (e *Engine) writeResults(ctx context.Context, tableCfg config.Table, columnNames []string, generators map[string]generator.Generator, outputChan <-chan rowData, tx db.Tx, singleUpdateQuery string, results *TableReport, p producer.Producer) ([]Diff, error) {
	var diffs []Diff
	var bar *progressbar.ProgressBar
	if !e.DryRun {
		count, _ := p.EstimateCount(ctx)
		bar = progressbar.Default(count, fmt.Sprintf("Masking %s", tableCfg.Name))
	}

	buffer := make([]rowData, 0, e.BatchSize)

	flush := func() error {
		if len(buffer) == 0 {
			return nil
		}
		if err := e.applyBatch(ctx, tx, tableCfg, columnNames, generators, buffer, singleUpdateQuery, results); err != nil {
			return err
		}
		if bar != nil {
			_ = bar.Add(len(buffer))
		}
		buffer = buffer[:0]
		return nil
	}

	for row := range outputChan {
		if e.DryRun {
			results.Updated++
			for i, colName := range columnNames {
				diffs = append(diffs, Diff{
					TableName:  tableCfg.Name,
					ColumnName: colName,
					PKValue:    row.pkValue,
					OldValue:   row.oldValues[i],
					NewValue:   row.newValues[i],
				})
			}
		} else {
			buffer = append(buffer, row)
			if len(buffer) >= e.BatchSize {
				if err := flush(); err != nil {
					return nil, err
				}
			}
		}
	}

	if err := flush(); err != nil {
		return nil, err
	}

	if bar != nil {
		_ = bar.Finish()
		fmt.Println()
	}

	return diffs, nil
}

func (e *Engine) applyBatch(ctx context.Context, tx db.Tx, tableCfg config.Table, columnNames []string, generators map[string]generator.Generator, buffer []rowData, singleUpdateQuery string, results *TableReport) error {
	query := e.buildBatchUpdateQuery(tableCfg.Name, tableCfg.PK, columnNames, results.Types, len(buffer))
	args := make([]any, 0, len(buffer)*(len(columnNames)+1))
	for _, row := range buffer {
		args = append(args, row.pkValue)
		args = append(args, row.newValues...)
	}

	// Use a savepoint to protect the transaction from batch failure
	_ = tx.Exec(ctx, "SAVEPOINT poof_batch")

	err := tx.Exec(ctx, query, args...)
	if err == nil {
		results.Updated += int64(len(buffer))
		_ = tx.Exec(ctx, "RELEASE SAVEPOINT poof_batch")
		return nil
	}

	// Rollback to savepoint so we can continue with individual updates
	_ = tx.Exec(ctx, "ROLLBACK TO SAVEPOINT poof_batch")

	// On batch failure, fallback to row-by-row
	slog.Warn("batch update failed, falling back to individual updates", "error", err, "table", tableCfg.Name, "batch_size", len(buffer))
	for _, row := range buffer {
		r := row
		if err := e.retryUpdate(ctx, tx, singleUpdateQuery, tableCfg, columnNames, generators, &r, maxRetries, results); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) retryUpdate(ctx context.Context, tx db.Tx, query string, tableCfg config.Table, columnNames []string, generators map[string]generator.Generator, row *rowData, retries int, results *TableReport) error {
	attempt := 0
	for {
		args := append([]any{row.pkValue}, row.newValues...)
		err := tx.Exec(ctx, query, args...)
		if err == nil {
			results.Updated++
			return nil
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && attempt < retries {
			attempt++
			results.Retried++
			// Re-generate values for retry
			for j, colName := range columnNames {
				seedBy := tableCfg.Columns[j].SeedBy
				newVal, genErr := generators[colName].Generate(generator.NewRowContext(
					tableCfg.Name,
					colName,
					e.Config.Locale,
					e.Config.Safety.Salt,
					row.pkValue,
					row.oldValues[j],
					seedBy,
					attempt, // Uses current attempt from unique_violation loop
				))
				if genErr != nil {
					return fmt.Errorf("retry gen error: %w", genErr)
				}
				row.newValues[j] = newVal
			}
			continue
		}

		results.Failed++
		slog.Error("row masking failed",
			"table", tableCfg.Name,
			"pk", row.pkValue,
			"error", err,
		)
		return nil // Non-fatal
	}
}
