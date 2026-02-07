package producer

import "sync"

var registerOnce sync.Once

func RegisterAll() {
	registerOnce.Do(func() {
		Register("table", NewTableProducer)
		Register("view", NewViewProducer)
		Register("query", NewQueryProducer)
	})
}
