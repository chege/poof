package db

import (
	"context"
)

// Rows abstracts driver-specific row iteration.
type Rows interface {
	Next() bool
	Values() ([]any, error)
	Scan(dest ...any) error
	Close()
}

// Tx abstracts database transactions.
type Tx interface {
	Exec(ctx context.Context, sql string, arguments ...any) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// DB defines the minimal interface for a database backend.
type DB interface {
	GetDatabaseName(ctx context.Context) (string, error)
	EstimateRowCount(ctx context.Context, tableName string) (int64, error)
	FetchRows(ctx context.Context, tableName string, pkColumn string, columns []string, limit int) (Rows, error)
	Begin(ctx context.Context) (Tx, error)
	Close()
}
