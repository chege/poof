package postgres

import (
	"context"
	"fmt"

	"github.com/chege/poof/internal/db"
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
			) as is_pk,
			EXISTS (
				SELECT 1 
				FROM information_schema.key_column_usage kcu
				JOIN information_schema.table_constraints tc ON tc.constraint_name = kcu.constraint_name
				WHERE kcu.table_name = c.table_name 
				  AND kcu.column_name = c.column_name 
				  AND tc.constraint_type = 'FOREIGN KEY'
			) as is_fk
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
		if err := rows.Scan(&col.Name, &col.DataType, &col.IsNullable, &col.HasUnique, &col.IsPrimaryKey, &col.IsForeignKey); err != nil {
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

// GetForeignKeys returns all foreign key relationships where the given table is the source.
func (c *Client) GetForeignKeys(ctx context.Context, tableName string) ([]db.ForeignKey, error) {
	query := `
		SELECT
			kcu.table_name,
			kcu.column_name,
			ccu.table_name AS referenced_table_name,
			ccu.column_name AS referenced_column_name
		FROM information_schema.table_constraints AS tc
		JOIN information_schema.key_column_usage AS kcu
		  ON tc.constraint_name = kcu.constraint_name
		  AND tc.table_schema = kcu.table_schema
		JOIN information_schema.constraint_column_usage AS ccu
		  ON ccu.constraint_name = tc.constraint_name
		WHERE tc.constraint_type = 'FOREIGN KEY' 
		  AND kcu.table_name = $1;
	`

	rows, err := c.pool.Query(ctx, query, tableName)
	if err != nil {
		return nil, fmt.Errorf("querying foreign keys: %w", err)
	}
	defer rows.Close()

	var fks []db.ForeignKey
	for rows.Next() {
		var fk db.ForeignKey
		if err := rows.Scan(&fk.TableName, &fk.ColumnName, &fk.ReferencedTableName, &fk.ReferencedColumnName); err != nil {
			return nil, fmt.Errorf("scanning foreign key: %w", err)
		}
		fks = append(fks, fk)
	}

	return fks, nil
}
