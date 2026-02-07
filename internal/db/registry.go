// Package db provides database abstractions and backend-agnostic interfaces for dbmask.
package db

import (
	"context"
	"fmt"
	"strings"
)

// Factory is a function that creates a new DB instance from a connection string.
type Factory func(ctx context.Context, connStr string) (DB, error)

var factories = make(map[string]Factory)

// Register adds a new database factory for the given URI scheme (e.g. "postgres").
func Register(scheme string, factory Factory) {
	factories[scheme] = factory
}

// Connect parses the DSN scheme and instantiates the appropriate database backend.
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
