package postgres

import (
	"context"
	"fmt"
)

// SetJobState updates the global masking state in the database.
func (c *Client) SetJobState(ctx context.Context, status string) error {
	// 1. Ensure table exists
	_, err := c.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS _poof_state (
			id INT PRIMARY KEY,
			status TEXT,
			updated_at TIMESTAMP DEFAULT NOW()
		);`)
	if err != nil {
		return fmt.Errorf("failed to create state table: %w", err)
	}

	// 2. Upsert status
	_, err = c.pool.Exec(ctx, `
		INSERT INTO _poof_state (id, status, updated_at) 
		VALUES (1, $1, NOW())
		ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status, updated_at = NOW();
	`, status)
	if err != nil {
		return fmt.Errorf("failed to set job state: %w", err)
	}
	return nil
}

// GetJobState retrieves the current masking state from the database.
func (c *Client) GetJobState(ctx context.Context) (string, error) {
	var status string
	err := c.pool.QueryRow(ctx, "SELECT status FROM _poof_state WHERE id = 1").Scan(&status)
	if err != nil {
		// If table doesn't exist or row missing, we assume no job has ever run
		return "NONE", nil
	}
	return status, nil
}
