// Package db provides database abstractions and backend-agnostic interfaces for poof.
package db

import (
	"context"
)

// Rows abstracts driver-specific row iteration.
type Rows interface {
	Next() bool
	Values() ([]any, error)
	Scan(dest ...any) error
	Err() error
	Close() error
}

// Tx abstracts database transactions.
type Tx interface {
	Exec(ctx context.Context, sql string, arguments ...any) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// ColumnInfo represents metadata about a database column.
type ColumnInfo struct {
	Name         string
	DataType     string
	IsNullable   bool
	HasUnique    bool
	IsPrimaryKey bool
	IsForeignKey bool
}

// ForeignKey represents a relational link between two tables.
type ForeignKey struct {
	TableName            string
	ColumnName           string
	ReferencedTableName  string
	ReferencedColumnName string
}

// DB defines the minimal interface for a database backend.
type DB interface {
	GetDatabaseName(ctx context.Context) (string, error)
	GetTableColumns(ctx context.Context, tableName string) ([]ColumnInfo, error)
	GetForeignKeys(ctx context.Context, tableName string) ([]ForeignKey, error)
	GetAllTables(ctx context.Context) ([]string, error)
	GetJobState(ctx context.Context) (string, error)
	SetJobState(ctx context.Context, status string) error
	EstimateRowCount(ctx context.Context, tableName string) (int64, error)
	FetchRows(ctx context.Context, tableName string, pkColumn string, columns []string, filter string, limit int) (Rows, error)
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	Begin(ctx context.Context) (Tx, error)
	Close()
}
