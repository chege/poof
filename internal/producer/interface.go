package producer

import (
	"context"
	"github.com/christopher/masker/internal/db"
)

type Producer interface {
	EstimateCount(ctx context.Context) (int64, error)
	FetchRows(ctx context.Context, columns []string, limit int) (db.Rows, error)
}
