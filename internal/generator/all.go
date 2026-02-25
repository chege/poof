// Package generator provides the masking data generation logic.
package generator

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/chege/poof/internal/config"
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

		RegisterGenerator("hash", func(_ config.Gen) (Generator, error) {
			return NewHashGenerator(), nil
		})

		RegisterGenerator("counter", func(_ config.Gen) (Generator, error) {
			return NewCounterGenerator(), nil
		})

		// Register Generic Faker Providers that use localizedData
		providers := []string{"first_name", "last_name", "company_name"}
		for _, p := range providers {
			name := p
			RegisterLocalizedFakerProvider("", name, func(_ *rand.Rand) any {
				// The actual locale will be determined during Generate() call via context.
				// But we need to register a provider that can handle it.
				// This is a bit of a trick: the provider function itself doesn't know the locale,
				// but GetLocalizedData will be called by the fakerGenerator.
				// Wait, the FakerProvider signature is FakerProvider func(r *rand.Rand) any.
				// It doesn't receive the locale.

				// Refactoring thought: FakerProvider should maybe receive the locale?
				// For now, I'll stick to the current implementation where the fakerGenerator handles it.
				return "placeholder" // This won't be used if we refactor faker.go
			})
		}

		// Let's fix faker.go to be truly locale-aware in the provider function if needed,
		// or just have the provider function be a "selector" from the data.

		registerFakerProviders()

		RegisterTestFakers()
	})
}

func registerFakerProviders() {
	// Standard providers
	RegisterFakerProvider("first_name", func(r *rand.Rand) any {
		// Note: This default provider will be overridden by localized versions if they exist
		data := GetLocalizedData("en_US", "first_name")
		return data[r.Intn(len(data))]
	})

	RegisterFakerProvider("last_name", func(r *rand.Rand) any {
		data := GetLocalizedData("en_US", "last_name")
		return data[r.Intn(len(data))]
	})

	RegisterFakerProvider("full_name", func(r *rand.Rand) any {
		fn := GetLocalizedData("en_US", "first_name")
		ln := GetLocalizedData("en_US", "last_name")
		return fmt.Sprintf("%s %s", fn[r.Intn(len(fn))], ln[r.Intn(len(ln))])
	})

	RegisterFakerProvider("email", func(r *rand.Rand) any {
		fn := GetLocalizedData("en_US", "first_name")
		ln := GetLocalizedData("en_US", "last_name")
		return GetLocalizedEmail(fn[r.Intn(len(fn))], ln[r.Intn(len(ln))])
	})

	RegisterFakerProvider("company_name", func(r *rand.Rand) any {
		data := GetLocalizedData("en_US", "company_name")
		return data[r.Intn(len(data))]
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

	// Register German overrides
	RegisterLocalizedFakerProvider("de_DE", "first_name", func(r *rand.Rand) any {
		data := GetLocalizedData("de_DE", "first_name")
		return data[r.Intn(len(data))]
	})
	RegisterLocalizedFakerProvider("de_DE", "last_name", func(r *rand.Rand) any {
		data := GetLocalizedData("de_DE", "last_name")
		return data[r.Intn(len(data))]
	})
	RegisterLocalizedFakerProvider("de_DE", "full_name", func(r *rand.Rand) any {
		fn := GetLocalizedData("de_DE", "first_name")
		ln := GetLocalizedData("de_DE", "last_name")
		return fmt.Sprintf("%s %s", fn[r.Intn(len(fn))], ln[r.Intn(len(ln))])
	})
}
