package producer

import (
	"context"
	"fmt"

	"github.com/christopher/poof/internal/config"
	"github.com/christopher/poof/internal/db"
)

type viewProducer struct {
	database db.DB
	viewName string
	pk       string
}

// NewViewProducer creates a new Producer that selects rows from a database view.
func NewViewProducer(_ context.Context, database db.DB, _ string, pk string, cfg *config.Source) (Producer, error) {
	if cfg == nil || cfg.Name == "" {
		return nil, fmt.Errorf("view producer requires a 'name'")
	}
	return &viewProducer{database: database, viewName: cfg.Name, pk: pk}, nil
}

func (p *viewProducer) EstimateCount(ctx context.Context) (int64, error) {
	return p.database.EstimateRowCount(ctx, p.viewName)
}

func (p *viewProducer) FetchRows(ctx context.Context, columns []string, filter string, limit int) (db.Rows, error) {
	return p.database.FetchRows(ctx, p.viewName, p.pk, columns, filter, limit)
}
