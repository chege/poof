// Package generator provides the masking data generation logic.
package generator

import (
	"math/rand"
)

// RegisterTestFakers registers a set of predictable faker providers for use in automated tests.
func RegisterTestFakers() {
	RegisterFakerProvider("test_name", func(_ *rand.Rand) any {
		return "Test User"
	})

	RegisterFakerProvider("test_email", func(_ *rand.Rand) any {
		return "test@example.com"
	})

	RegisterFakerProvider("test_username", func(_ *rand.Rand) any {
		return "test_user_1"
	})

	RegisterFakerProvider("test_company", func(_ *rand.Rand) any {
		return "Test Corp"
	})

	RegisterFakerProvider("test_phone", func(_ *rand.Rand) any {
		return "+1-555-0000"
	})

	RegisterFakerProvider("test_ipv4", func(_ *rand.Rand) any {
		return "127.0.0.1"
	})

	RegisterFakerProvider("test_short_text", func(_ *rand.Rand) any {
		return "test_text"
	})
}
