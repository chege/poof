package masker

import (
	"context"
	"testing"
	"time"

	"github.com/christopher/masker/internal/config"
	"github.com/christopher/masker/internal/db"
	"github.com/christopher/masker/internal/generator"
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
	client, err := db.NewClient(ctx, connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	_, err = client.Pool.Exec(ctx, `
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT
		);
		INSERT INTO users (name, email) VALUES ('Real Name 1', 'real1@example.com');
		INSERT INTO users (name, email) VALUES ('Real Name 2', 'real2@example.com');
	`)
	if err != nil {
		t.Fatal(err)
	}

	// 3. Mask
	generator.RegisterAll()
	cfg := &config.Config{
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
	rows, _ := client.Pool.Query(ctx, "SELECT name, email FROM users ORDER BY id")
	defer rows.Close()

	for rows.Next() {
		var name, email string
		rows.Scan(&name, &email)
		assert.Equal(t, "Test User", name)
		assert.Equal(t, "test@example.com", email)
	}
}

func TestE2E_NewProviders(t *testing.T) {
	ctx := context.Background()

	pgContainer, _ := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("pass"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2)),
	)
	defer pgContainer.Terminate(ctx)
	connStr, _ := pgContainer.ConnectionString(ctx, "sslmode=disable")
	client, _ := db.NewClient(ctx, connStr)
	defer client.Close()

	client.Pool.Exec(ctx, `
		CREATE TABLE profiles (
			id SERIAL PRIMARY KEY,
			username TEXT,
			company TEXT,
			phone TEXT,
			ip TEXT,
			bio TEXT
		);
		INSERT INTO profiles (username, company, phone, ip, bio) VALUES ('real_user', 'Real Corp', '123', '1.2.3.4', 'real bio');
	`)

	generator.RegisterAll()
	cfg := &config.Config{
		Tables: []config.Table{{
			Name: "profiles", PK: "id",
			Columns: []config.Column{
				{Name: "username", Gen: config.Gen{Type: "faker", Provider: "test_username"}},
				{Name: "company", Gen: config.Gen{Type: "faker", Provider: "test_company"}},
				{Name: "phone", Gen: config.Gen{Type: "faker", Provider: "test_phone"}},
				{Name: "ip", Gen: config.Gen{Type: "faker", Provider: "test_ipv4"}},
				{Name: "bio", Gen: config.Gen{Type: "faker", Provider: "test_short_text"}},
			},
		}},
	}

	engine := NewEngine(client, cfg, 1)
	_, err := engine.Apply(ctx)
	assert.NoError(t, err)

	var username, company, phone, ip, bio string
	err = client.Pool.QueryRow(ctx, "SELECT username, company, phone, ip, bio FROM profiles WHERE id = 1").Scan(&username, &company, &phone, &ip, &bio)
	assert.NoError(t, err)
	assert.Equal(t, "test_user_1", username)
	assert.Equal(t, "Test Corp", company)
	assert.Equal(t, "+1-555-0000", phone)
	assert.Equal(t, "127.0.0.1", ip)
	assert.Equal(t, "test_text", bio)
}

func TestDeterminism(t *testing.T) {
	ctx := context.Background()

	pgContainer, _ := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("pass"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2)),
	)
	defer pgContainer.Terminate(ctx)
	connStr, _ := pgContainer.ConnectionString(ctx, "sslmode=disable")
	client, _ := db.NewClient(ctx, connStr)
	defer client.Close()

	client.Pool.Exec(ctx, "CREATE TABLE users (id INT PRIMARY KEY, name TEXT)")
	client.Pool.Exec(ctx, "INSERT INTO users (id, name) VALUES (1, 'Real'), (2, 'Real'), (3, 'Real')")

	generator.RegisterAll()
	cfg := &config.Config{
		Tables: []config.Table{{
			Name: "users", PK: "id",
			Columns: []config.Column{{Name: "name", Gen: config.Gen{Type: "faker", Provider: "first_name"}}},
		}},
	}

	// Run with 1 worker
	engine1 := NewEngine(client, cfg, 1)
	engine1.Apply(ctx)
	rows1, _ := client.Pool.Query(ctx, "SELECT id, name FROM users ORDER BY id")
	results1 := make(map[int]string)
	for rows1.Next() {
		var id int
		var name string
		rows1.Scan(&id, &name)
		results1[id] = name
	}
	rows1.Close()

	// Reset and run with 4 workers
	client.Pool.Exec(ctx, "UPDATE users SET name = 'Real'")
	engine4 := NewEngine(client, cfg, 4)
	engine4.Apply(ctx)
	rows4, _ := client.Pool.Query(ctx, "SELECT id, name FROM users ORDER BY id")
	for rows4.Next() {
		var id int
		var name string
		rows4.Scan(&id, &name)
		assert.Equal(t, results1[id], name, "Output should be deterministic regardless of workers")
	}
}
