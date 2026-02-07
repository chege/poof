// Package producer defines the interfaces and implementations for selecting rows from a database.
package producer

import (
	"context"
	"github.com/christopher/masker/internal/db"
)

// Producer is the interface for identifying which rows should be masked.
type Producer interface {
	// EstimateCount returns an approximate number of rows that will be produced.
	EstimateCount(ctx context.Context) (int64, error)
	// FetchRows returns a cursor for iterating over the selected rows.
	FetchRows(ctx context.Context, columns []string, limit int) (db.Rows, error)
}
