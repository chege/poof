package generator

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type counterGenerator struct {
	counts sync.Map
}

// NewCounterGenerator creates a new generator that returns incrementing integers.
func NewCounterGenerator() Generator {
	return &counterGenerator{}
}

func (g *counterGenerator) Generate(ctx RowContext) (any, error) {
	key := fmt.Sprintf("%s:%s", ctx.TableName, ctx.ColumnName)
	val, _ := g.counts.LoadOrStore(key, new(int64))
	ptr, ok := val.(*int64)
	if !ok {
		return nil, fmt.Errorf("unexpected type in counter storage")
	}
	return atomic.AddInt64(ptr, 1), nil
}
