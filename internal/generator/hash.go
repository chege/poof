package generator

import (
	"fmt"
)

type hashGenerator struct{}

// NewHashGenerator creates a new generator that produces deterministic MD5 hashes.
func NewHashGenerator() Generator {
	return &hashGenerator{}
}

func (g *hashGenerator) Generate(ctx RowContext) (any, error) {
	return fmt.Sprintf("%x", ctx.Seed), nil
}
