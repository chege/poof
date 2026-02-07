package producer

import (
	"context"
	"fmt"
	"strings"

	"github.com/christopher/masker/internal/config"
	"github.com/christopher/masker/internal/db"
)

type queryProducer struct {
	database db.DB
	sql      string
	pk       string
}

func NewQueryProducer(ctx context.Context, database db.DB, table string, pk string, cfg *config.Source) (Producer, error) {
	if cfg == nil || cfg.SQL == "" {
		return nil, fmt.Errorf("query producer requires 'sql'")
	}

	sql := strings.TrimSpace(cfg.SQL)
	upperSQL := strings.ToUpper(sql)

	if !strings.HasPrefix(upperSQL, "SELECT") {
		return nil, fmt.Errorf("query must be a SELECT statement")
	}

	if !strings.Contains(upperSQL, "ORDER BY") {
		return nil, fmt.Errorf("query must include ORDER BY to ensure determinism")
	}

	return &queryProducer{database: database, sql: sql, pk: pk}, nil
}

func (p *queryProducer) EstimateCount(ctx context.Context) (int64, error) {
	// For queries, we return 0 as a default estimate.
	return 0, nil
}

func (p *queryProducer) FetchRows(ctx context.Context, columns []string, limit int) (db.Rows, error) {
	sql := p.sql
	if limit > 0 {
		sql = fmt.Sprintf("SELECT * FROM (%s) AS subquery LIMIT %d", p.sql, limit)
	}
	return p.database.Query(ctx, sql)
}
