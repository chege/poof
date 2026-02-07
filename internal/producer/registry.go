package producer

import (
	"context"
	"fmt"
	"sync"

	"github.com/christopher/masker/internal/config"
	"github.com/christopher/masker/internal/db"
)

type Factory func(ctx context.Context, database db.DB, table string, pk string, cfg *config.Source) (Producer, error)

var (
	registryMu sync.RWMutex
	factories  = make(map[string]Factory)
)

func Register(name string, factory Factory) {
	registryMu.Lock()
	defer registryMu.Unlock()
	if _, dup := factories[name]; dup {
		panic(fmt.Sprintf("producer factory %q already registered", name))
	}
	factories[name] = factory
}

func NewProducer(ctx context.Context, database db.DB, table string, pk string, cfg *config.Source) (Producer, error) {
	registryMu.RLock()
	defer registryMu.RUnlock()

	sourceType := "table"
	if cfg != nil {
		sourceType = cfg.Type
	}

	factory, ok := factories[sourceType]
	if !ok {
		return nil, fmt.Errorf("producer type %q not found", sourceType)
	}
	return factory(ctx, database, table, pk, cfg)
}
