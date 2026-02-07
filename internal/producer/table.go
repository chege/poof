package producer

import (
	"context"
	"github.com/christopher/masker/internal/config"
	"github.com/christopher/masker/internal/db"
)

type tableProducer struct {
	database db.DB
	table    string
	pk       string
}

// NewTableProducer creates a new Producer that performs a standard table scan.
func NewTableProducer(_ context.Context, database db.DB, table string, pk string, _ *config.Source) (Producer, error) {
	return &tableProducer{database: database, table: table, pk: pk}, nil
}

func (p *tableProducer) EstimateCount(ctx context.Context) (int64, error) {
	return p.database.EstimateRowCount(ctx, p.table)
}

func (p *tableProducer) FetchRows(ctx context.Context, columns []string, limit int) (db.Rows, error) {
	return p.database.FetchRows(ctx, p.table, p.pk, columns, limit)
}
