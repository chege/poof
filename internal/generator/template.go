package generator

import (
	"bytes"
	"text/template"
)

type templateGenerator struct {
	tmpl *template.Template
}

func NewTemplateGenerator(text string) (Generator, error) {
	tmpl, err := template.New("gen").Parse(text)
	if err != nil {
		return nil, err
	}
	return &templateGenerator{tmpl: tmpl}, nil
}

func (g *templateGenerator) Generate(ctx RowContext) (any, error) {
	var buf bytes.Buffer
	err := g.tmpl.Execute(&buf, ctx)
	if err != nil {
		return nil, err
	}
	return buf.String(), nil
}
