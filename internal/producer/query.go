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

// NewQueryProducer creates a new Producer based on a custom SQL query.
func NewQueryProducer(_ context.Context, database db.DB, _ string, pk string, cfg *config.Source) (Producer, error) {
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

// EstimateCount returns an estimate of the number of rows (currently not implemented for custom queries).
func (p *queryProducer) EstimateCount(_ context.Context) (int64, error) {
	// For queries, we return 0 as a default estimate.
	return 0, nil
}

// FetchRows executes the custom SQL query and returns an iterator over the result rows.
func (p *queryProducer) FetchRows(ctx context.Context, _ []string, limit int) (db.Rows, error) {
	sql := p.sql
	if limit > 0 {
		sql = fmt.Sprintf("SELECT * FROM (%s) AS subquery LIMIT %d", p.sql, limit)
	}
	return p.database.Query(ctx, sql)
}
