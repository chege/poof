package generator

import (
	"fmt"
	"sync"

	"github.com/christopher/masker/internal/config"
)

type GeneratorFactory func(gen config.Gen) (Generator, error)

var (
	registryMu sync.RWMutex
	factories  = make(map[string]GeneratorFactory)
)

func RegisterGenerator(name string, factory GeneratorFactory) {
	registryMu.Lock()
	defer registryMu.Unlock()
	if _, dup := factories[name]; dup {
		panic(fmt.Sprintf("generator %q already registered", name))
	}
	factories[name] = factory
}

func NewGenerator(gen config.Gen) (Generator, error) {
	registryMu.RLock()
	factory, ok := factories[gen.Type]
	registryMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("generator %q not found", gen.Type)
	}
	return factory(gen)
}