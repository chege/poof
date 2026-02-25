// Package postgres provides the PostgreSQL backend implementation for poof.
package postgres

import (
	"context"

	"github.com/chege/poof/internal/db"
)

// init registers the PostgreSQL backend with the database registry.
func init() {
	db.Register("postgres", func(ctx context.Context, connStr string) (db.DB, error) {
		return NewClient(ctx, connStr)
	})
}
