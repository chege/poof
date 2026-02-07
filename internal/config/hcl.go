package config

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

func LoadConfig(path string) (*Config, error) {
	var cfg Config
	err := hclsimple.DecodeFile(path, nil, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if len(cfg.Tables) == 0 {
		return nil, fmt.Errorf("config must contain at least one table block")
	}

	return &cfg, nil
}

func (c *Config) IsAllowed(dbName string) bool {
	if len(c.Allowlist) == 0 {
		return false
	}
	for _, allowed := range c.Allowlist {
		if allowed == dbName {
			return true
		}
	}
	return false
}
