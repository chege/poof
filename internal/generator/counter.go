package generator

import (
	"encoding/binary"
)

type counterGenerator struct {
}

// NewCounterGenerator creates a new generator that returns deterministic integers based on the row seed.
// Note: This replaces the global incrementing counter to maintain parallelism-safe determinism.
func NewCounterGenerator() Generator {
	return &counterGenerator{}
}

func (g *counterGenerator) Generate(ctx RowContext) (any, error) {
	// Use the first 8 bytes of the deterministic seed as an int64
	// #nosec G115 -- seed conversion is intentional.
	val := int64(binary.BigEndian.Uint64(ctx.Seed[:8]))
	if val < 0 {
		val = -val
	}
	return val, nil
}

func (g *counterGenerator) ExpectedType() string {
	return "int64"
}
