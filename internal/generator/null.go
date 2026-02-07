// Package generator provides the masking data generation logic.
package generator

type nullGenerator struct{}

// NewNullGenerator creates a new generator that always returns nil (SQL NULL).
func NewNullGenerator() Generator {
	return &nullGenerator{}
}

// Generate returns nil, ignoring the row context.
func (g *nullGenerator) Generate(_ RowContext) (any, error) {
	return nil, nil
}
