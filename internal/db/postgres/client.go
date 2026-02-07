package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/christopher/masker/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	pool *pgxpool.Pool
}

func NewClient(ctx context.Context, connStr string) (*Client, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to postgres: %w", err)
	}
	return &Client{pool: pool}, nil
}

func (c *Client) Close() {
	c.pool.Close()
}

func (c *Client) GetDatabaseName(ctx context.Context) (string, error) {
	var dbName string
	err := c.pool.QueryRow(ctx, "SELECT current_database()").Scan(&dbName)
	if err != nil {
		return "", fmt.Errorf("failed to get database name: %w", err)
	}
	return dbName, nil
}

func (c *Client) EstimateRowCount(ctx context.Context, tableName string) (int64, error) {
	var count int64
	// Fast estimate for Postgres
	err := c.pool.QueryRow(ctx, "SELECT reltuples::bigint FROM pg_class WHERE relname = $1", tableName).Scan(&count)
	if err != nil || count < 0 {
		// Fallback to COUNT(*) if estimate fails or returns -1 (not yet analyzed)
		err = c.pool.QueryRow(ctx, fmt.Sprintf("SELECT count(*) FROM %s", tableName)).Scan(&count)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *Client) FetchRows(ctx context.Context, tableName string, pkColumn string, columns []string, limit int) (db.Rows, error) {
	cols := []string{pkColumn}
	cols = append(cols, columns...)

	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY %s", strings.Join(cols, ", "), tableName, pkColumn)
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := c.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return &rowsWrapper{rows}, nil
}

func (c *Client) Begin(ctx context.Context) (db.Tx, error) {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &txWrapper{tx}, nil
}

func (c *Client) Query(ctx context.Context, sql string, args ...any) (db.Rows, error) {
	rows, err := c.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &rowsWrapper{rows}, nil
}

type rowsWrapper struct {
	pgx.Rows
}

func (w *rowsWrapper) Values() ([]any, error) {
	return w.Rows.Values()
}

func (w *rowsWrapper) Err() error {
	return w.Rows.Err()
}

func (w *rowsWrapper) Close() error {
	w.Rows.Close()
	return nil
}

type txWrapper struct {
	pgx.Tx
}

func (w *txWrapper) Exec(ctx context.Context, sql string, arguments ...any) error {
	_, err := w.Tx.Exec(ctx, sql, arguments...)
	return err
}
