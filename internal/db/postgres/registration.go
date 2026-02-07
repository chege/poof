// Package postgres provides the PostgreSQL backend implementation for dbmask.
package postgres

import (
	"context"

	"github.com/christopher/masker/internal/db"
)

// init registers the PostgreSQL backend with the database registry.
func init() {
	db.Register("postgres", func(ctx context.Context, connStr string) (db.DB, error) {
		return NewClient(ctx, connStr)
	})
}
