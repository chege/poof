// Package generator provides the masking data generation logic.
package generator

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/christopher/masker/internal/config"
)

var registerOnce sync.Once

// RegisterAll ensures all built-in generators and faker providers are registered exactly once.
func RegisterAll() {
	registerOnce.Do(func() {
		RegisterGenerator("faker", func(gen config.Gen) (Generator, error) {
			if gen.Provider == "" {
				return nil, fmt.Errorf("faker generator requires a provider")
			}
			return NewFakerGenerator(gen.Provider), nil
		})

		RegisterGenerator("constant", func(gen config.Gen) (Generator, error) {
			return NewConstantGenerator(gen.Value), nil
		})

		RegisterGenerator("null", func(_ config.Gen) (Generator, error) {
			return NewNullGenerator(), nil
		})

		RegisterGenerator("template", func(gen config.Gen) (Generator, error) {
			if gen.Template == "" {
				return nil, fmt.Errorf("template generator requires a template")
			}
			return NewTemplateGenerator(gen.Template)
		})

		// Register Faker Providers
		RegisterFakerProvider("first_name", func(r *rand.Rand) any {
			names := []string{"John", "Jane", "Alice", "Bob", "Charlie", "Diana"}
			return names[r.Intn(len(names))]
		})
		RegisterFakerProvider("last_name", func(r *rand.Rand) any {
			names := []string{"Doe", "Smith", "Johnson", "Williams", "Brown", "Jones"}
			return names[r.Intn(len(names))]
		})
		RegisterFakerProvider("email", func(r *rand.Rand) any {
			first := []string{"john", "jane", "alice", "bob"}
			last := []string{"doe", "smith", "johnson"}
			domains := []string{"example.com", "test.org", "mail.com"}
			return fmt.Sprintf("%s.%s@%s", first[r.Intn(len(first))], last[r.Intn(len(last))], domains[r.Intn(len(domains))])
		})

		RegisterFakerProvider("username", func(r *rand.Rand) any {
			first := []string{"john", "jane", "alice", "bob"}
			last := []string{"doe", "smith", "johnson"}
			return fmt.Sprintf("%s_%s%d", first[r.Intn(len(first))], last[r.Intn(len(last))], r.Intn(100))
		})

		RegisterFakerProvider("company_name", func(r *rand.Rand) any {
			names := []string{"Acme Corp", "Globex", "Soylent Corp", "Initech", "Umbrella Corp"}
			return names[r.Intn(len(names))]
		})

		RegisterFakerProvider("phone_number", func(r *rand.Rand) any {
			return fmt.Sprintf("+1-555-%04d", r.Intn(10000))
		})

		RegisterFakerProvider("ipv4_address", func(r *rand.Rand) any {
			return fmt.Sprintf("%d.%d.%d.%d", r.Intn(256), r.Intn(256), r.Intn(256), r.Intn(256))
		})

		RegisterFakerProvider("short_text", func(r *rand.Rand) any {
			chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
			b := make([]byte, 10)
			for i := range b {
				b[i] = chars[r.Intn(len(chars))]
			}
			return string(b)
		})

		RegisterTestFakers()
	})
}
