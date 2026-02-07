// Package generator provides the masking data generation logic.
package generator

type constantGenerator struct {
	value string
}

// NewConstantGenerator creates a new generator that always returns the same fixed value.
func NewConstantGenerator(value string) Generator {
	return &constantGenerator{value: value}
}

// Generate returns the constant value, ignoring the row context.
func (g *constantGenerator) Generate(_ RowContext) (any, error) {
	return g.value, nil
}
