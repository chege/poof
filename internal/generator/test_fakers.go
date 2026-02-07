package generator

import (
	"math/rand"
)

func RegisterTestFakers() {
	RegisterFakerProvider("test_name", func(r *rand.Rand) any {
		return "Test User"
	})

	RegisterFakerProvider("test_email", func(r *rand.Rand) any {
		return "test@example.com"
	})

	RegisterFakerProvider("test_username", func(r *rand.Rand) any {
		return "test_user_1"
	})

	RegisterFakerProvider("test_company", func(r *rand.Rand) any {
		return "Test Corp"
	})

	RegisterFakerProvider("test_phone", func(r *rand.Rand) any {
		return "+1-555-0000"
	})

	RegisterFakerProvider("test_ipv4", func(r *rand.Rand) any {
		return "127.0.0.1"
	})

	RegisterFakerProvider("test_short_text", func(r *rand.Rand) any {
		return "test_text"
	})
}
