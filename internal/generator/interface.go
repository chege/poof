// Package generator provides the masking data generation logic.
package generator

// RowContext carries information about the current row and column being masked.
type RowContext struct {
	PrimaryKeyValue any
	OriginalValue   any
	TableName       string
	ColumnName      string
	Locale          string
	Salt            string
	Attempt         int      // For retrying generation
	Seed            [16]byte // MD5 hash used for deterministic generation
}

// Generator is the interface for producing masked values.
type Generator interface {
	// Generate produces a new value based on the row context.
	Generate(ctx RowContext) (any, error)
	// ExpectedType returns the Go type string this generator produces (e.g. "string", "int64").
	ExpectedType() string
}
