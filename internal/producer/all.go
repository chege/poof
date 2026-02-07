// Package producer defines the interfaces and implementations for selecting rows from a database.
package producer

import "sync"

var registerOnce sync.Once

// RegisterAll ensures all built-in producers are registered exactly once.
func RegisterAll() {
	registerOnce.Do(func() {
		Register("table", NewTableProducer)
		Register("view", NewViewProducer)
		Register("query", NewQueryProducer)
	})
}
