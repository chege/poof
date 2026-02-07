package masker

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/christopher/masker/internal/config"
	"github.com/christopher/masker/internal/db"
	"github.com/christopher/masker/internal/generator"
)

type Engine struct {
	DB      *db.Client
	Config  *config.Config
	Workers int
}

func NewEngine(client *db.Client, cfg *config.Config, workers int) *Engine {
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
	pkValue any
	values  []any
}

func (e *Engine) Apply(ctx context.Context) error {
	for _, tableCfg := range e.Config.Tables {
		if err := e.maskTable(ctx, tableCfg); err != nil {
			return fmt.Errorf("failed to mask table %s: %w", tableCfg.Name, err)
		}
	}
	return nil
}

func (e *Engine) maskTable(ctx context.Context, tableCfg config.Table) error {
	log.Printf("Masking table %s with %d workers...", tableCfg.Name, e.Workers)

	columnNames := make([]string, 0, len(tableCfg.Columns))
	generators := make(map[string]generator.Generator)

	for _, col := range tableCfg.Columns {
		columnNames = append(columnNames, col.Name)
		gen, err := generator.NewGenerator(col.Gen)
		if err != nil {
			return fmt.Errorf("failed to create generator for column %s: %w", col.Name, err)
		}
		generators[col.Name] = gen
	}

	rows, err := e.DB.FetchRows(ctx, tableCfg.Name, tableCfg.PK, columnNames)
	if err != nil {
		return fmt.Errorf("failed to fetch rows: %w", err)
	}
	defer rows.Close()

	tx, err := e.DB.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	updateQuery := fmt.Sprintf("UPDATE %s SET ", tableCfg.Name)
	for i, col := range columnNames {
		if i > 0 {
			updateQuery += ", "
		}
		updateQuery += fmt.Sprintf("%s = $%d", col, i+2)
	}
	updateQuery += fmt.Sprintf(" WHERE %s = $1", tableCfg.PK)

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
				row.values = newValues
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
			inputChan <- rowData{pkValue: values[0]}
		}
	}()

	// Writer loop
	for row := range outputChan {
		select {
		case err := <-errChan:
			return err
		default:
		}

		args := append([]any{row.pkValue}, row.values...)
		_, err = tx.Exec(ctx, updateQuery, args...)
		if err != nil {
			return fmt.Errorf("failed to update row: %w", err)
		}
	}

	return tx.Commit(ctx)
}