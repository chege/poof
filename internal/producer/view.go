package producer

import (
	"context"
	"fmt"
	"github.com/christopher/masker/internal/config"
	"github.com/christopher/masker/internal/db"
)

type viewProducer struct {
	database db.DB
	viewName string
	pk       string
}

func NewViewProducer(ctx context.Context, database db.DB, table string, pk string, cfg *config.Source) (Producer, error) {
	if cfg == nil || cfg.Name == "" {
		return nil, fmt.Errorf("view producer requires a 'name'")
	}
	return &viewProducer{database: database, viewName: cfg.Name, pk: pk}, nil
}

func (p *viewProducer) EstimateCount(ctx context.Context) (int64, error) {
	return p.database.EstimateRowCount(ctx, p.viewName)
}

func (p *viewProducer) FetchRows(ctx context.Context, columns []string, limit int) (db.Rows, error) {
	return p.database.FetchRows(ctx, p.viewName, p.pk, columns, limit)
}
