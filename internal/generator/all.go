package generator

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/christopher/masker/internal/config"
)

var registerOnce sync.Once

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

	RegisterGenerator("null", func(gen config.Gen) (Generator, error) {
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

	RegisterFakerProvider("full_name", func(r *rand.Rand) any {
		firsts := []string{"John", "Jane", "Alice", "Bob"}
		lasts := []string{"Doe", "Smith", "Johnson", "Williams"}
		return fmt.Sprintf("%s %s", firsts[r.Intn(len(firsts))], lasts[r.Intn(len(lasts))])
	})

	RegisterFakerProvider("phone_number", func(r *rand.Rand) any {
		return fmt.Sprintf("+1-555-%04d", r.Intn(10000))
	})

	RegisterFakerProvider("company", func(r *rand.Rand) any {
		names := []string{"Acme Corp", "Globex", "Soylent Corp", "Initech", "Umbrella Corp"}
		return names[r.Intn(len(names))]
	})

	RegisterTestFakers()
	})
}
