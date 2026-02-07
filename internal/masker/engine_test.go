package masker

import (
	"context"
	"testing"
	"time"

	"github.com/christopher/masker/internal/config"
	"github.com/christopher/masker/internal/db"
	_ "github.com/christopher/masker/internal/db/postgres"
	"github.com/christopher/masker/internal/generator"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
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
	defer pgContainer.Terminate(ctx)

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
		Database: config.Database{DSN: connStr},
		Safety:   config.Safety{AllowedDBNames: []string{"testdb"}},
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
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, email string
		rows.Scan(&id, &name, &email)
		if diff := cmp.Diff("Test User", name); diff != "" {
			t.Errorf("name mismatch (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff("test@example.com", email); diff != "" {
			t.Errorf("email mismatch (-want +got):\n%s", diff)
		}
	}
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
	defer pgContainer.Terminate(ctx)
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
		Database: config.Database{DSN: connStr},
		Safety:   config.Safety{AllowedDBNames: []string{"testdb"}},
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
	defer rows.Close()
	rows.Next()
	var id int
	var name string
	rows.Scan(&id, &name)
	assert.Equal(t, "Real Name", name, "Database should be unchanged in dry-run")
}

func TestE2E_QueryProducer(t *testing.T) {
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
	defer pgContainer.Terminate(ctx)
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
	err = tx.Exec(ctx, `
		CREATE TABLE users (id INT PRIMARY KEY, name TEXT, active BOOLEAN);
		INSERT INTO users (id, name, active) VALUES (1, 'Real 1', true), (2, 'Real 2', false);
	`)
	assert.NoError(t, err)
	err = tx.Commit(ctx)
	assert.NoError(t, err)

	generator.RegisterAll()
	cfg := &config.Config{
		Database: config.Database{DSN: connStr},
		Safety:   config.Safety{AllowedDBNames: []string{"testdb"}},
		Tables: []config.Table{{
			Name: "users", PK: "id",
			Source: &config.Source{
				Type: "query",
				SQL:  "SELECT id FROM users WHERE active = true ORDER BY id",
			},
			Columns: []config.Column{{Name: "name", Gen: config.Gen{Type: "constant", Value: "Masked"}}},
		}},
	}

	engine := NewEngine(client, cfg, 1)
	_, err = engine.Apply(ctx)
	assert.NoError(t, err)

	// Verify only active user is masked
	var name1, name2 string

	rows1, err := client.Query(ctx, "SELECT name FROM users WHERE id = 1")
	assert.NoError(t, err)
	rows1.Next()
	rows1.Scan(&name1)
	rows1.Close()

	rows2, err := client.Query(ctx, "SELECT name FROM users WHERE id = 2")
	assert.NoError(t, err)
	rows2.Next()
	rows2.Scan(&name2)
	rows2.Close()

	assert.Equal(t, "Masked", name1)
	assert.Equal(t, "Real 2", name2)
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
	defer pgContainer.Terminate(ctx)
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
		Database: config.Database{DSN: connStr},
		Safety:   config.Safety{AllowedDBNames: []string{"testdb"}},
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
		rows1.Scan(&id, &name)
		results1[id] = name
	}
	rows1.Close()

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
		rows4.Scan(&id, &name)
		if diff := cmp.Diff(results1[id], name); diff != "" {
			t.Errorf("determinism mismatch for ID %d (-want +got):\n%s", id, diff)
		}
	}
	rows4.Close()
}
