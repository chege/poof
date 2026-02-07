package masker

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/christopher/masker/internal/config"
	"github.com/christopher/masker/internal/db"
	"github.com/christopher/masker/internal/generator"
	"github.com/christopher/masker/internal/producer"
)

type Engine struct {
	DB      db.DB
	Config  *config.Config
	Workers int
	DryRun  bool
}

type TableReport struct {
	Name      string
	Estimates int64
	Columns   []string
}

type MaskingReport struct {
	Tables []TableReport
	Diffs  []Diff
}

type Diff struct {
	TableName  string
	ColumnName string
	PKValue    any
	OldValue   any
	NewValue   any
}

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
			count = 0 // Ignore estimate errors
		}

		cols := make([]string, len(tableCfg.Columns))
		for i, c := range tableCfg.Columns {
			cols[i] = c.Name
		}

		report.Tables = append(report.Tables, TableReport{
			Name:      tableCfg.Name,
			Estimates: count,
			Columns:   cols,
		})

		diffs, err := e.maskTable(ctx, tableCfg, p)
		if err != nil {
			return nil, fmt.Errorf("failed to mask table %s: %w", tableCfg.Name, err)
		}
		report.Diffs = append(report.Diffs, diffs...)
	}
	return report, nil
}

func (e *Engine) maskTable(ctx context.Context, tableCfg config.Table, p producer.Producer) ([]Diff, error) {
	if !e.DryRun {
		log.Printf("Masking table %s with %d workers...", tableCfg.Name, e.Workers)
	}

	columnNames := make([]string, 0, len(tableCfg.Columns))
	generators := make(map[string]generator.Generator)

	for _, col := range tableCfg.Columns {
		columnNames = append(columnNames, col.Name)
		gen, err := generator.NewGenerator(col.Gen)
		if err != nil {
			return nil, fmt.Errorf("failed to create generator for column %s: %w", col.Name, err)
		}
		generators[col.Name] = gen
	}

	// Limit fetching in DryRun for plan preview
	limit := 0
	if e.DryRun {
		limit = 5
	}

	rows, err := p.FetchRows(ctx, columnNames, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rows: %w", err)
	}
	defer rows.Close()

	var tx db.Tx
	var updateQuery string
	if !e.DryRun {
		var err error
		tx, err = e.DB.Begin(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to start transaction: %w", err)
		}
		defer tx.Rollback(ctx)

		updateQuery = fmt.Sprintf("UPDATE %s SET ", tableCfg.Name)
		for i, col := range columnNames {
			if i > 0 {
				updateQuery += ", "
			}
			updateQuery += fmt.Sprintf("%s = $%d", col, i+2)
		}
		updateQuery += fmt.Sprintf(" WHERE %s = $1", tableCfg.PK)
	}

	inputChan := make(chan rowData, e.Workers*2)
	outputChan := make(chan rowData, e.Workers*2)
	errChan := make(chan error, 1)

	var wg sync.WaitGroup
	for i := 0; i < e.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range inputChan {
				newValues := make([]any, len(columnNames))
				for j, colName := range columnNames {
					gen := generators[colName]
					rowCtx := generator.NewRowContext(tableCfg.Name, colName, row.pkValue)
					val, err := gen.Generate(rowCtx)
					if err != nil {
						select {
						case errChan <- err:
						default:
						}
						return
					}
					newValues[j] = val
				}
				row.newValues = newValues
				outputChan <- row
			}
		}()
	}

	go func() {
		wg.Wait()
		close(outputChan)
	}()

	// Reader goroutine
	go func() {
		defer close(inputChan)
		for rows.Next() {
			values, err := rows.Values()
			if err != nil {
				select {
				case errChan <- err:
				default:
				}
				return
			}
			inputChan <- rowData{pkValue: values[0], oldValues: values[1:]}
		}
	}()

	var diffs []Diff
	// Writer/Collector loop
	for row := range outputChan {
		select {
		case err := <-errChan:
			return nil, err
		default:
		}

		if e.DryRun {
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
			args := append([]any{row.pkValue}, row.newValues...)
			err = tx.Exec(ctx, updateQuery, args...)
			if err != nil {
				return nil, fmt.Errorf("failed to update row: %w", err)
			}
		}
	}

	if !e.DryRun {
		return nil, tx.Commit(ctx)
	}
	return diffs, nil
}
