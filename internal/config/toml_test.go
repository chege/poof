package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_EnvExpansion(t *testing.T) {
	// Create a temporary config file with env vars
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.toml")

	content := `
[databases.staging]
dsn = "postgres://user:${DB_PASSWORD}@localhost:5432/db"

[safety]
allowed_db_names = ["db"]

[[tables]]
name = "users"
pk = "id"
  [[tables.columns]]
  name = "email"
  [tables.columns.gen]
  type = "faker"
  provider = "email"
`
	err := os.WriteFile(configPath, []byte(content), 0600)
	assert.NoError(t, err)

	// Set the env var
	err = os.Setenv("DB_PASSWORD", "secret123")
	assert.NoError(t, err)
	defer func() {
		_ = os.Unsetenv("DB_PASSWORD")
	}()

	// Load the config
	cfg, err := LoadConfig(configPath)
	assert.NoError(t, err)

	// Verify expansion
	assert.Equal(t, "postgres://user:secret123@localhost:5432/db", cfg.Databases["staging"].DSN)
}
