// Package generator provides the masking data generation logic.
package generator

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/christopher/poof/internal/config"
)

// Factory is a function that creates a new Generator based on its configuration.
type Factory func(gen config.Gen) (Generator, error)

var (
	registryMu sync.RWMutex
	factories  = make(map[string]Factory)
)

// RegisterGenerator adds a new generator factory to the global registry.
func RegisterGenerator(name string, factory Factory) {
	registryMu.Lock()
	defer registryMu.Unlock()
	if _, dup := factories[name]; dup {
		panic(fmt.Sprintf("generator %q already registered", name))
	}
	slog.Debug("Registering generator", "name", name)
	factories[name] = factory
}

// NewGenerator instantiates a generator based on the column configuration.
func NewGenerator(gen config.Gen) (Generator, error) {
	registryMu.RLock()
	factory, ok := factories[gen.Type]
	registryMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("generator %q not found", gen.Type)
	}
	return factory(gen)
}

// ValidateGen checks if the generator configuration is semantically valid.
func ValidateGen(locale string, gen config.Gen) error {
	registryMu.RLock()
	_, ok := factories[gen.Type]
	registryMu.RUnlock()

	if !ok {
		return fmt.Errorf("unknown generator type %q", gen.Type)
	}

	switch gen.Type {
	case "faker":
		if !ProviderExists(locale, gen.Provider) {
			return fmt.Errorf("unknown faker provider %q", gen.Provider)
		}
	case "template":
		if err := ValidateTemplate(gen.Template); err != nil {
			return fmt.Errorf("invalid template: %w", err)
		}
	}

	return nil
}

// Validator implements the config.GeneratorValidator interface.
type Validator struct{}

// ValidateGen checks if the generator configuration is semantically valid.
func (v *Validator) ValidateGen(locale string, gen config.Gen) error {
	return ValidateGen(locale, gen)
}
