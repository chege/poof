// Package generator provides the masking data generation logic.
package generator

import (
	"bytes"
	"text/template"
)

type templateGenerator struct {
	tmpl *template.Template
}

// NewTemplateGenerator creates a new generator that uses Go text/template for value generation.
func NewTemplateGenerator(text string) (Generator, error) {
	tmpl, err := template.New("gen").Parse(text)
	if err != nil {
		return nil, err
	}
	return &templateGenerator{tmpl: tmpl}, nil
}

// Generate executes the template with the given row context.
func (g *templateGenerator) Generate(ctx RowContext) (any, error) {
	var buf bytes.Buffer
	err := g.tmpl.Execute(&buf, ctx)
	if err != nil {
		return nil, err
	}
	return buf.String(), nil
}
