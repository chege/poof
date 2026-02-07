package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	Pool *pgxpool.Pool
}

func NewClient(ctx context.Context, connStr string) (*Client, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	return &Client{Pool: pool}, nil
}

func (c *Client) Close() {
	c.Pool.Close()
}

func (c *Client) GetDatabaseName(ctx context.Context) (string, error) {
	var dbName string
	err := c.Pool.QueryRow(ctx, "SELECT current_database()").Scan(&dbName)
	if err != nil {
		return "", fmt.Errorf("failed to get database name: %w", err)
	}
	return dbName, nil
}

func (c *Client) FetchRows(ctx context.Context, tableName string, pkColumn string, columns []string) (pgx.Rows, error) {
	cols := pkColumn
	for _, col := range columns {
		cols += ", " + col
	}
	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY %s", cols, tableName, pkColumn)
	return c.Pool.Query(ctx, query)
}
