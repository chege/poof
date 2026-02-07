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

	// Strict validation
	if err := validateBody(file.Body, configSchema); err != nil {
		return nil, err
	}

	syntaxBody := file.Body.(*hclsyntax.Body)
	for _, block := range syntaxBody.Blocks {
		if block.Type != "table" {
			continue
		}
		if err := validateBody(block.Body, tableSchema); err != nil {
			return nil, fmt.Errorf("table %q: %w", block.Labels[0], err)
		}

		for _, subBlock := range block.Body.Blocks {
			if subBlock.Type == "column" {
				if err := validateBody(subBlock.Body, columnSchema); err != nil {
					return nil, fmt.Errorf("table %q column %q: %w", block.Labels[0], subBlock.Labels[0], err)
				}
				for _, gen := range subBlock.Body.Blocks {
					if gen.Type == "gen" {
						if err := validateBody(gen.Body, genSchema); err != nil {
							return nil, fmt.Errorf("table %q column %q gen %q: %w", block.Labels[0], subBlock.Labels[0], gen.Labels[0], err)
						}
					}
				}
			} else if subBlock.Type == "source" {
				if err := validateBody(subBlock.Body, sourceSchema); err != nil {
					return nil, fmt.Errorf("table %q source %q: %w", block.Labels[0], subBlock.Labels[0], err)
				}
			}
		}
	}

	return &cfg, nil
}

func validateBody(body hcl.Body, schema *hcl.BodySchema) error {
	content, remainingBody, diags := body.PartialContent(schema)
	if diags.HasErrors() {
		return diags
	}

	// Check for unknown items in the remaining body
	content, _, _ = remainingBody.PartialContent(&hcl.BodySchema{})
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
		{Type: "table", LabelNames: []string{"name"}},
	},
}

var tableSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "pk", Required: true},
	},
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "column", LabelNames: []string{"name"}},
		{Type: "source", LabelNames: []string{"type"}},
	},
}

var columnSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "gen", LabelNames: []string{"type"}},
	},
}

var genSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "provider", Required: false},
		{Name: "value", Required: false},
		{Name: "template", Required: false},
		{Name: "params", Required: false},
	},
}

var sourceSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "name", Required: false},
		{Name: "sql", Required: false},
		{Name: "params", Required: false},
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
