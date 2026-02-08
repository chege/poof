package postgres

import (
	"context"
	"fmt"

	"github.com/christopher/poof/internal/db"
)

// GetTableColumns returns metadata for all columns in a given table.
func (c *Client) GetTableColumns(ctx context.Context, tableName string) ([]db.ColumnInfo, error) {
	query := `
		SELECT 
			column_name, 
			data_type, 
			is_nullable = 'YES' as is_nullable,
			EXISTS (
				SELECT 1 
				FROM information_schema.constraint_column_usage cu
				JOIN information_schema.table_constraints tc ON tc.constraint_name = cu.constraint_name
				WHERE cu.table_name = c.table_name 
				  AND cu.column_name = c.column_name 
				  AND tc.constraint_type IN ('UNIQUE', 'PRIMARY KEY')
			) as has_unique,
			EXISTS (
				SELECT 1 
				FROM information_schema.constraint_column_usage cu
				JOIN information_schema.table_constraints tc ON tc.constraint_name = cu.constraint_name
				WHERE cu.table_name = c.table_name 
				  AND cu.column_name = c.column_name 
				  AND tc.constraint_type = 'PRIMARY KEY'
			) as is_pk
		FROM information_schema.columns c
		WHERE table_name = $1
		ORDER BY ordinal_position;
	`

	rows, err := c.pool.Query(ctx, query, tableName)
	if err != nil {
		return nil, fmt.Errorf("querying column info: %w", err)
	}
	defer rows.Close()

	var columns []db.ColumnInfo
	for rows.Next() {
		var col db.ColumnInfo
		if err := rows.Scan(&col.Name, &col.DataType, &col.IsNullable, &col.HasUnique, &col.IsPrimaryKey); err != nil {
			return nil, fmt.Errorf("scanning column info: %w", err)
		}
		columns = append(columns, col)
	}

	return columns, nil
}

// GetAllTables returns a list of all user tables in the current database.
func (c *Client) GetAllTables(ctx context.Context) ([]string, error) {
	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		  AND table_type = 'BASE TABLE'
		ORDER BY table_name;
	`

	rows, err := c.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scanning table name: %w", err)
		}
		tables = append(tables, name)
	}

	return tables, nil
}
