// Package generator provides the masking data generation logic.
package generator

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/chege/poof/internal/config"
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
func ValidateGen(locale string, gen config.Gen, sqlDataType string) error {
	registryMu.RLock()
	factory, ok := factories[gen.Type]
	registryMu.RUnlock()

	if !ok {
		return fmt.Errorf("unknown generator type %q", gen.Type)
	}

	g, err := factory(gen)
	if err != nil {
		return err
	}

	// 1. Static parameter checks
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

	// 2. Type compatibility checks (if SQL type is provided)
	if sqlDataType != "" {
		producedType := g.ExpectedType()
		if !isTypeCompatible(producedType, sqlDataType) {
			return fmt.Errorf("type mismatch: generator produces %q but column is %q", producedType, sqlDataType)
		}
	}

	return nil
}

func isTypeCompatible(goType, sqlType string) bool {
	sqlType = strings.ToLower(sqlType)
	switch goType {
	case "string":
		// Strings can go into almost anything via stringification
		return true
	case "int64":
		return strings.Contains(sqlType, "int") || strings.Contains(sqlType, "serial") || strings.Contains(sqlType, "numeric") || strings.Contains(sqlType, "text") || strings.Contains(sqlType, "char")
	case "bool":
		return sqlType == "boolean" || strings.Contains(sqlType, "text") || strings.Contains(sqlType, "char")
	case "null":
		return true
	}
	return true
}

// Validator implements the config.GeneratorValidator interface.
type Validator struct{}

// ValidateGen checks if the generator configuration is semantically valid.
func (v *Validator) ValidateGen(locale string, gen config.Gen, sqlDataType string) error {
	return ValidateGen(locale, gen, sqlDataType)
}
