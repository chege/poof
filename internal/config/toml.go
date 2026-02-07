// Package config handles the loading, validation, and representation of the poof TOML configuration.
package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// LoadConfig reads and validates a TOML configuration file from the given path.
func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to open config: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()

	var cfg Config
	decoder := toml.NewDecoder(f)
	// Strict validation: check for unknown keys
	metadata, err := decoder.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TOML: %w", err)
	}

	if len(metadata.Undecoded()) > 0 {
		return nil, fmt.Errorf("config contains unknown fields: %v", metadata.Undecoded())
	}

	// Run declarative validation
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// Save writes the configuration to the given writer in TOML format.
func (c *Config) Save(w io.Writer) error {
	return toml.NewEncoder(w).Encode(c)
}
