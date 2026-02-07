// Package engine coordinates the data masking process.
package engine

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgconn"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/sync/errgroup"

	"github.com/christopher/poof/internal/config"
	"github.com/christopher/poof/internal/db"
	"github.com/christopher/poof/internal/generator"
	"github.com/christopher/poof/internal/producer"
)

// Engine orchestrates the data masking process across multiple tables and workers.
type Engine struct {
	DB       db.DB
	Config   *config.Config
	Workers  int
	DryRun   bool
	PlanOnly bool
}

// TableReport provides a summary of masking results for a specific table.
type TableReport struct {
	Name      string
	Columns   []string
	Estimates int64
	Updated   int64
	Retried   int64
	Failed    int64
}

// MaskingReport provides a comprehensive summary of masking operations.
type MaskingReport struct {
	Tables []TableReport
	Diffs  []Diff
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
	return &Engine{
		DB:      client,
		Config:  cfg,
		Workers: workers,
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
	report := &MaskingReport{}
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
	return report, nil
}

func (e *Engine) maskTable(ctx context.Context, tableCfg config.Table, p producer.Producer) ([]Diff, TableReport, error) {
	results := TableReport{Name: tableCfg.Name}
	columnNames := make([]string, 0, len(tableCfg.Columns))
	generators := make(map[string]generator.Generator)

	for _, col := range tableCfg.Columns {
		columnNames = append(columnNames, col.Name)
		gen, err := generator.NewGenerator(col.Gen)
		if err != nil {
			return nil, results, fmt.Errorf("failed to create generator for column %s: %w", col.Name, err)
		}
		generators[col.Name] = gen
	}

	limit := 0
	if e.DryRun {
		limit = 5
	}

	rows, err := p.FetchRows(ctx, columnNames, limit)
	if err != nil {
		return nil, results, fmt.Errorf("failed to fetch rows: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var tx db.Tx
	var updateQuery string
	if !e.DryRun {
		var err error
		tx, err = e.DB.Begin(ctx)
		if err != nil {
			return nil, results, fmt.Errorf("failed to start transaction: %w", err)
		}
		defer func() {
			if tx != nil {
				_ = tx.Rollback(ctx)
			}
		}()

		updateQuery = fmt.Sprintf("UPDATE %s SET ", tableCfg.Name)
		for i, col := range columnNames {
			if i > 0 {
				updateQuery += ", "
			}
			updateQuery += fmt.Sprintf("%s = $%d", col, i+2)
		}
		updateQuery += fmt.Sprintf(" WHERE %s = $1", tableCfg.PK)
	}

	g, gCtx := errgroup.WithContext(ctx)
	inputChan := make(chan rowData, e.Workers*2)
	outputChan := make(chan rowData, e.Workers*2)

	// Reader
	g.Go(func() error {
		defer close(inputChan)
		for rows.Next() {
			values, err := rows.Values()
			if err != nil {
				return fmt.Errorf("failed to read values: %w", err)
			}
			select {
			case inputChan <- rowData{pkValue: values[0], oldValues: values[1:]}:
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
					newValues := make([]any, len(columnNames))
					for j, colName := range columnNames {
						val, err := generators[colName].Generate(generator.NewRowContext(tableCfg.Name, colName, row.pkValue))
						if err != nil {
							return fmt.Errorf("gen error: %w", err)
						}
						newValues[j] = val
					}
					row.newValues = newValues
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

	var diffs []Diff
	var bar *progressbar.ProgressBar
	if !e.DryRun {
		count, _ := p.EstimateCount(ctx)
		bar = progressbar.Default(count, fmt.Sprintf("Masking %s", tableCfg.Name))
	}

	maxRetries := 10

	// Writer loop
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
			attempt := 0
			for {
				args := append([]any{row.pkValue}, row.newValues...)
				err := tx.Exec(ctx, updateQuery, args...)
				if err == nil {
					results.Updated++
					break
				}

				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) && pgErr.Code == "23505" && attempt < maxRetries {
					attempt++
					results.Retried++
					// Re-generate values for retry
					for j, colName := range columnNames {
						newVal, genErr := generators[colName].Generate(generator.NewRowContext(tableCfg.Name, colName, row.pkValue))
						if genErr != nil {
							return nil, results, fmt.Errorf("retry gen error: %w", genErr)
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
				break
			}
			if bar != nil {
				_ = bar.Add(1)
			}
		}
	}

	if err := g.Wait(); err != nil {
		return nil, results, fmt.Errorf("table %s masking failed: %w", tableCfg.Name, err)
	}

	if !e.DryRun {
		if bar != nil {
			_ = bar.Finish()
			fmt.Println()
		}
		if err := tx.Commit(ctx); err != nil {
			return nil, results, fmt.Errorf("commit error: %w", err)
		}
		tx = nil
	}

	return diffs, results, nil
}
