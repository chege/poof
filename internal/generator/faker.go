// Package generator provides the masking data generation logic.
package generator

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
)

// FakerProvider is a function that generates a fake value using the provided random source.
type FakerProvider func(r *rand.Rand) any

var (
	fakerMu sync.RWMutex
	// locale -> providerName -> provider
	fakerProviders = make(map[string]map[string]FakerProvider)
)

// DefaultLocale is the fallback locale used when none is specified.
const DefaultLocale = "en_US"

// RegisterFakerProvider adds a new faker provider to the registry for the default locale.
func RegisterFakerProvider(name string, provider FakerProvider) {
	RegisterLocalizedFakerProvider(DefaultLocale, name, provider)
}

// RegisterLocalizedFakerProvider adds a new faker provider for a specific locale.
func RegisterLocalizedFakerProvider(locale, name string, provider FakerProvider) {
	fakerMu.Lock()
	defer fakerMu.Unlock()
	if fakerProviders[locale] == nil {
		fakerProviders[locale] = make(map[string]FakerProvider)
	}
	slog.Debug("Registering faker provider", "locale", locale, "name", name)
	fakerProviders[locale][name] = provider
}

// GetFakerProvider retrieves a faker provider by name and locale from the registry.
func GetFakerProvider(locale, name string) (FakerProvider, bool) {
	fakerMu.RLock()
	defer fakerMu.RUnlock()

	// Try requested locale
	if providers, ok := fakerProviders[locale]; ok {
		if p, ok := providers[name]; ok {
			return p, true
		}
	}

	// Fallback to default locale
	if providers, ok := fakerProviders[DefaultLocale]; ok {
		if p, ok := providers[name]; ok {
			return p, true
		}
	}

	return nil, false
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
	locale := ctx.Locale
	if locale == "" {
		locale = DefaultLocale
	}

	provider, ok := GetFakerProvider(locale, g.providerName)
	if !ok {
		return nil, fmt.Errorf("faker provider %q not found for locale %q", g.providerName, locale)
	}

	// Create a local rand source from the seed for determinism.
	// We use the first 8 bytes of the MD5 hash.
	// #nosec G115 -- seed conversion is intentional and non-critical.
	seedInt := int64(binary.BigEndian.Uint64(ctx.Seed[:8]))

	// #nosec G404 -- math/rand is used for deterministic seeding, not for cryptographic security.
	r := rand.New(rand.NewSource(seedInt))

	return provider(r), nil
}
