package generator

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"sync"
)

type FakerProvider func(r *rand.Rand) any

var (
	fakerMu        sync.RWMutex
	fakerProviders = make(map[string]FakerProvider)
)

func RegisterFakerProvider(name string, provider FakerProvider) {
	fakerMu.Lock()
	defer fakerMu.Unlock()
	fakerProviders[name] = provider
}

func GetFakerProvider(name string) (FakerProvider, bool) {
	fakerMu.RLock()
	defer fakerMu.RUnlock()
	p, ok := fakerProviders[name]
	return p, ok
}

type fakerGenerator struct {
	providerName string
}

func NewFakerGenerator(providerName string) Generator {
	return &fakerGenerator{providerName: providerName}
}

func (g *fakerGenerator) Generate(ctx RowContext) (any, error) {
	fakerMu.RLock()
	provider, ok := fakerProviders[g.providerName]
	fakerMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("faker provider %q not found", g.providerName)
	}

	// Create a local rand source from the seed
	seedInt := int64(binary.BigEndian.Uint64(ctx.Seed[:8]))
	r := rand.New(rand.NewSource(seedInt))

	return provider(r), nil
}
