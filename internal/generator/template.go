// Package generator provides the masking data generation logic.
package generator

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/christopher/poof/internal/config"
)

// ValidateTemplate checks if a template string is syntactically correct.
func ValidateTemplate(text string) error {
	// We use the same dummy FuncMap as NewTemplateGenerator
	funcMap := template.FuncMap{
		"faker":   func(_ string) (string, error) { return "", nil },
		"counter": func() (string, error) { return "", nil },
		"hash":    func() (string, error) { return "", nil },
	}
	_, err := template.New("validate").Funcs(funcMap).Parse(text)
	return err
}

type templateGenerator struct {
	tmpl *template.Template
}

// NewTemplateGenerator creates a new generator that uses Go text/template for value generation.
func NewTemplateGenerator(text string) (Generator, error) {
	tg := &templateGenerator{}

	// We use a dummy context during parsing because FuncMap must be defined upfront.
	// The actual data will be provided during Execute.
	funcMap := template.FuncMap{
		"faker":   func(_ string) (string, error) { return "", nil },
		"counter": func() (string, error) { return "", nil },
		"hash":    func() (string, error) { return "", nil },
	}

	tmpl, err := template.New("gen").Funcs(funcMap).Parse(text)
	if err != nil {
		return nil, fmt.Errorf("template parse error: %w", err)
	}
	tg.tmpl = tmpl

	return tg, nil
}

// Generate executes the template with the given row context.
func (g *templateGenerator) Generate(ctx RowContext) (any, error) {
	// We re-inject the functions with the current RowContext closure
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

	// Create a clone of the template to safely override the FuncMap for this execution
	// without affecting other concurrent workers using the same generator instance.
	t, err := g.tmpl.Clone()
	if err != nil {
		return nil, err
	}
	t.Funcs(funcMap)

	var buf bytes.Buffer
	err = t.Execute(&buf, ctx)
	if err != nil {
		return nil, fmt.Errorf("template execution error: %w", err)
	}
	return buf.String(), nil
}

func (g *templateGenerator) ExpectedType() string {
	return "string"
}
