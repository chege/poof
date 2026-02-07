package generator

import (
	"math/rand"
)

func RegisterTestFakers() {
	RegisterFakerProvider("test_name", func(r *rand.Rand) any {
		// We use a fixed sequence or something predictable.
		// Since we seed the rand with the row seed, we can just use Intn.
		// But for "test fakers", we might want something even more obvious.
		return "Test User"
	})

	RegisterFakerProvider("test_email", func(r *rand.Rand) any {
		return "test@example.com"
	})
}
