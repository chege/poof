// Package generator provides the masking data generation logic.
package generator

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/christopher/poof/internal/config"
)

type templateGenerator struct {
	text string
}

// NewTemplateGenerator creates a new generator that uses Go text/template for value generation.
func NewTemplateGenerator(text string) (Generator, error) {
	// We only store the text here because we need to inject the FuncMap
	// which requires access to the RowContext during execution.
	return &templateGenerator{text: text}, nil
}

// Generate executes the template with the given row context.
func (g *templateGenerator) Generate(ctx RowContext) (any, error) {
	// Define helper functions for the template
	funcMap := template.FuncMap{
		"faker": func(provider string) (string, error) {
			gen, err := NewGenerator(config.Gen{
				Type:     "faker",
				Provider: provider,
			})
			if err != nil {
				return "", err
			}
			val, err := gen.Generate(ctx)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%v", val), nil
		},
		"counter": func() (string, error) {
			gen, err := NewGenerator(config.Gen{Type: "counter"})
			if err != nil {
				return "", err
			}
			val, err := gen.Generate(ctx)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%v", val), nil
		},
		"hash": func() (string, error) {
			gen, err := NewGenerator(config.Gen{Type: "hash"})
			if err != nil {
				return "", err
			}
			val, err := gen.Generate(ctx)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%v", val), nil
		},
	}

	tmpl, err := template.New("gen").Funcs(funcMap).Parse(g.text)
	if err != nil {
		return nil, fmt.Errorf("template parse error: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, ctx)
	if err != nil {
		return nil, fmt.Errorf("template execution error: %w", err)
	}
	return buf.String(), nil
}
