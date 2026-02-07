package config

import (
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config: %w", err)
	}
	defer f.Close()

	src, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	file, diags := hclsyntax.ParseConfig(src, path, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse config: %s", diags.Error())
	}

	var cfg Config
	decodeDiags := gohcl.DecodeBody(file.Body, nil, &cfg)
	if decodeDiags.HasErrors() {
		return nil, fmt.Errorf("failed to decode config: %s", decodeDiags.Error())
	}

	// Strict validation for the root body
	_, remainingBody, diags := file.Body.PartialContent(configSchema)
	if diags.HasErrors() {
		return nil, fmt.Errorf("root validation error: %s", diags.Error())
	}
	if err := checkUnknown(remainingBody); err != nil {
		return nil, fmt.Errorf("root config error: %w", err)
	}

	return &cfg, nil
}

func checkUnknown(body hcl.Body) error {
	content, _, _ := body.PartialContent(&hcl.BodySchema{})
	for _, attr := range content.Attributes {
		return fmt.Errorf("unknown attribute %q at %s", attr.Name, attr.NameRange.String())
	}
	for _, block := range content.Blocks {
		return fmt.Errorf("unknown block %q at %s", block.Type, block.TypeRange.String())
	}
	return nil
}

var configSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "allowlist", Required: false},
	},
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "table",
			LabelNames: []string{"name"},
		},
	},
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
