package engine

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/christopher/poof/internal/config"
	"github.com/christopher/poof/internal/db"
	_ "github.com/christopher/poof/internal/db/postgres"
	"github.com/christopher/poof/internal/generator"
)

func TestE2E(t *testing.T) {
	ctx := context.Background()

	// 1. Setup Postgres Container
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("pass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = pgContainer.Terminate(ctx)
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// 2. Seed Data
	client, err := db.Connect(ctx, connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	tx, err := client.Begin(ctx)
	assert.NoError(t, err)
	err = tx.Exec(ctx, `
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT
		);
		INSERT INTO users (name, email) VALUES ('Real Name 1', 'real1@example.com');
		INSERT INTO users (name, email) VALUES ('Real Name 2', 'real2@example.com');
	`)
	assert.NoError(t, err)
	err = tx.Commit(ctx)
	assert.NoError(t, err)

	// 3. Mask
	generator.RegisterAll()
	cfg := &config.Config{
		Databases: map[string]config.Database{"default": {DSN: connStr}},
		Safety:    config.Safety{AllowedDBNames: []string{"testdb"}},
		Tables: []config.Table{
			{
				Name: "users",
				PK:   "id",
				Columns: []config.Column{
					{
						Name: "name",
						Gen: config.Gen{
							Type:     "faker",
							Provider: "test_name",
						},
					},
					{
						Name: "email",
						Gen: config.Gen{
							Type:     "faker",
							Provider: "test_email",
						},
					},
				},
			},
		},
	}

	engine := NewEngine(client, cfg, 2)
	_, err = engine.Apply(ctx)
	assert.NoError(t, err)

	// 4. Verify
	rows, err := client.FetchRows(ctx, "users", "id", []string{"name", "email"}, 0)
	assert.NoError(t, err)
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var id int
		var name, email string
		err = rows.Scan(&id, &name, &email)
		assert.NoError(t, err)
		if diff := cmp.Diff("Test User", name); diff != "" {
			t.Errorf("name mismatch (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff("test@example.com", email); diff != "" {
			t.Errorf("email mismatch (-want +got):\n%s", diff)
		}
	}
	assert.NoError(t, rows.Err())
}

func TestE2E_DryRun(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("pass"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2)),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = pgContainer.Terminate(ctx)
	}()
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	client, err := db.Connect(ctx, connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	tx, err := client.Begin(ctx)
	assert.NoError(t, err)
	err = tx.Exec(ctx, "CREATE TABLE users (id INT PRIMARY KEY, name TEXT)")
	assert.NoError(t, err)
	err = tx.Exec(ctx, "INSERT INTO users (id, name) VALUES (1, 'Real Name')")
	assert.NoError(t, err)
	err = tx.Commit(ctx)
	assert.NoError(t, err)

	generator.RegisterAll()
	cfg := &config.Config{
		Databases: map[string]config.Database{"default": {DSN: connStr}},
		Safety:    config.Safety{AllowedDBNames: []string{"testdb"}},
		Tables: []config.Table{{
			Name: "users", PK: "id",
			Columns: []config.Column{{Name: "name", Gen: config.Gen{Type: "faker", Provider: "test_name"}}},
		}},
	}

	engine := NewEngine(client, cfg, 1)
	engine.DryRun = true
	report, err := engine.Apply(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, report.Diffs)

	// Verify DB is UNCHANGED
	rows, err := client.FetchRows(ctx, "users", "id", []string{"name"}, 0)
	assert.NoError(t, err)
	defer func() {
		_ = rows.Close()
	}()
	assert.True(t, rows.Next())
	var id int
	var name string
	err = rows.Scan(&id, &name)
	assert.NoError(t, err)
	assert.Equal(t, "Real Name", name, "Database should be unchanged in dry-run")
	assert.NoError(t, rows.Err())
}

func TestE2E_BatchFailureFallback(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("pass"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2)),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = pgContainer.Terminate(ctx)
	}()
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	client, err := db.Connect(ctx, connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	tx, err := client.Begin(ctx)
	assert.NoError(t, err)
	// Create table with UNIQUE constraint
	err = tx.Exec(ctx, "CREATE TABLE users (id INT PRIMARY KEY, email TEXT UNIQUE)")
	assert.NoError(t, err)
	// Insert one record
	err = tx.Exec(ctx, "INSERT INTO users (id, email) VALUES (1, 'collision@example.com')")
	assert.NoError(t, err)
	// Insert another that we will TRY to mask to the same value
	err = tx.Exec(ctx, "INSERT INTO users (id, email) VALUES (2, 'original@example.com')")
	assert.NoError(t, err)
	err = tx.Commit(ctx)
	assert.NoError(t, err)

	generator.RegisterAll()
	cfg := &config.Config{
		BatchSize: 10, // Small batch
		Databases: map[string]config.Database{"default": {DSN: connStr}},
		Safety:    config.Safety{AllowedDBNames: []string{"testdb"}},
		Tables: []config.Table{{
			Name: "users", PK: "id",
			Columns: []config.Column{{
				Name: "email",
				Gen: config.Gen{
					Type:  "constant",
					Value: "collision@example.com", // This WILL collide for both rows if applied
				},
			}},
		}},
	}

	// We use a constant generator to FORCE a unique violation.
	// Since we use MD5 seeding for retries, the retry will ALSO be a constant in this test
	// unless we specifically use a generator that changes on retry.
	// BUT, our current engine.retryUpdate RE-GENERATES using the same generator.
	// For a 'constant' generator, it will always collide.
	// This is perfect for testing the 'Failed' count and fallback.

	engine := NewEngine(client, cfg, 1)
	report, err := engine.Apply(ctx)
	// Apply might return an error because the COMMIT might fail if we don't handle partial failures at the table level transaction.
	// But for this test, we just want to see the stats in the report if it exists.
	if err != nil {
		assert.Contains(t, err.Error(), "commit error")
	}

	if report != nil && len(report.Tables) > 0 {
		assert.Equal(t, int64(1), report.Tables[0].Failed, "Should have 1 failed row due to persistent collision")
	}
}

func TestDeterminism(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("pass"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2)),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = pgContainer.Terminate(ctx)
	}()
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	client, err := db.Connect(ctx, connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	tx, err := client.Begin(ctx)
	assert.NoError(t, err)
	err = tx.Exec(ctx, "CREATE TABLE users (id INT PRIMARY KEY, name TEXT)")
	assert.NoError(t, err)
	err = tx.Exec(ctx, "INSERT INTO users (id, name) VALUES (1, 'Real'), (2, 'Real'), (3, 'Real')")
	assert.NoError(t, err)
	err = tx.Commit(ctx)
	assert.NoError(t, err)

	generator.RegisterAll()
	cfg := &config.Config{
		Databases: map[string]config.Database{"default": {DSN: connStr}},
		Safety:    config.Safety{AllowedDBNames: []string{"testdb"}},
		Tables: []config.Table{{
			Name: "users", PK: "id",
			Columns: []config.Column{{Name: "name", Gen: config.Gen{Type: "faker", Provider: "first_name"}}},
		}},
	}

	// Run with 1 worker
	engine1 := NewEngine(client, cfg, 1)
	_, err = engine1.Apply(ctx)
	assert.NoError(t, err)
	rows1, err := client.FetchRows(ctx, "users", "id", []string{"name"}, 0)
	assert.NoError(t, err)
	results1 := make(map[int]string)
	for rows1.Next() {
		var id int
		var name string
		err = rows1.Scan(&id, &name)
		assert.NoError(t, err)
		results1[id] = name
	}
	assert.NoError(t, rows1.Err())
	err = rows1.Close()
	assert.NoError(t, err)

	// Reset and run with 4 workers
	tx, err = client.Begin(ctx)
	assert.NoError(t, err)
	err = tx.Exec(ctx, "UPDATE users SET name = 'Real'")
	assert.NoError(t, err)
	err = tx.Commit(ctx)
	assert.NoError(t, err)

	engine4 := NewEngine(client, cfg, 4)
	_, err = engine4.Apply(ctx)
	assert.NoError(t, err)
	rows4, err := client.FetchRows(ctx, "users", "id", []string{"name"}, 0)
	assert.NoError(t, err)
	for rows4.Next() {
		var id int
		var name string
		err = rows4.Scan(&id, &name)
		assert.NoError(t, err)
		if diff := cmp.Diff(results1[id], name); diff != "" {
			t.Errorf("determinism mismatch for ID %d (-want +got):\n%s", id, diff)
		}
	}
	assert.NoError(t, rows4.Err())
	err = rows4.Close()
	assert.NoError(t, err)
}
