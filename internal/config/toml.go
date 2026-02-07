package config

import (
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config: %w", err)
	}
	defer f.Close()

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

func (c *Config) Save(w io.Writer) error {
	return toml.NewEncoder(w).Encode(c)
}
