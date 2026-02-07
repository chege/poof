package db

import (
	"context"
	"fmt"
	"strings"
)

type Factory func(ctx context.Context, connStr string) (DB, error)

var factories = make(map[string]Factory)

func Register(scheme string, factory Factory) {
	factories[scheme] = factory
}

func Connect(ctx context.Context, connStr string) (DB, error) {
	parts := strings.Split(connStr, "://")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid connection string: missing scheme (e.g. postgres://)")
	}
	scheme := parts[0]

	factory, ok := factories[scheme]
	if !ok {
		return nil, fmt.Errorf("unsupported database scheme: %q", scheme)
	}

	return factory(ctx, connStr)
}
