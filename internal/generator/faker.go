// Package generator provides the masking data generation logic.
package generator

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"sync"
)

// FakerProvider is a function that generates a fake value using the provided random source.
type FakerProvider func(r *rand.Rand) any

var (
	fakerMu        sync.RWMutex
	fakerProviders = make(map[string]FakerProvider)
)

// RegisterFakerProvider adds a new faker provider to the registry.
func RegisterFakerProvider(name string, provider FakerProvider) {
	fakerMu.Lock()
	defer fakerMu.Unlock()
	fakerProviders[name] = provider
}

// GetFakerProvider retrieves a faker provider by name from the registry.
func GetFakerProvider(name string) (FakerProvider, bool) {
	fakerMu.RLock()
	defer fakerMu.RUnlock()
	p, ok := fakerProviders[name]
	return p, ok
}

type fakerGenerator struct {
	providerName string
}

// NewFakerGenerator creates a new generator that uses a faker provider for value generation.
func NewFakerGenerator(providerName string) Generator {
	return &fakerGenerator{providerName: providerName}
}

// Generate produces a fake value using the registered provider and a row-level deterministic seed.
func (g *fakerGenerator) Generate(ctx RowContext) (any, error) {
	fakerMu.RLock()
	provider, ok := fakerProviders[g.providerName]
	fakerMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("faker provider %q not found", g.providerName)
	}

	// Create a local rand source from the seed for determinism.
	// We use the first 8 bytes of the MD5 hash.
	// #nosec G115 -- seed conversion is intentional and non-critical.
	seedInt := int64(binary.BigEndian.Uint64(ctx.Seed[:8]))

	// #nosec G404 -- math/rand is used for deterministic seeding, not for cryptographic security.
	r := rand.New(rand.NewSource(seedInt))

	return provider(r), nil
}
