package postgres

import (
	"context"

	"github.com/christopher/masker/internal/db"
)

func init() {
	db.Register("postgres", func(ctx context.Context, connStr string) (db.DB, error) {
		return NewClient(ctx, connStr)
	})
}
