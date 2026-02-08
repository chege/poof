package generator

import "fmt"

// This file contains localized data for faker providers.
// In a future version, this could be loaded from YAML or JSON files.

var localizedData = map[string]map[string][]string{
	"en_US": {
		"first_name":   {"John", "Jane", "Alice", "Bob", "Charlie", "Diana"},
		"last_name":    {"Doe", "Smith", "Johnson", "Williams", "Brown", "Jones"},
		"company_name": {"Acme Corp", "Globex", "Soylent Corp", "Initech", "Umbrella Corp"},
	},
	"de_DE": {
		"first_name":   {"Hans", "Jürgen", "Karl", "Stefan", "Monika", "Angelika"},
		"last_name":    {"Müller", "Schmidt", "Schneider", "Fischer", "Weber", "Meyer"},
		"company_name": {"Müller GmbH", "Schmidt AG", "Schneider & Co"},
	},
}

// GetLocalizedData returns a slice of strings for a given locale and provider.
// It falls back to en_US if the locale is not found.
func GetLocalizedData(locale, provider string) []string {
	if data, ok := localizedData[locale]; ok {
		if slice, ok := data[provider]; ok {
			return slice
		}
	}
	// Fallback to en_US
	if data, ok := localizedData["en_US"]; ok {
		if slice, ok := data[provider]; ok {
			return slice
		}
	}
	return nil
}

// GetLocalizedEmail returns a deterministic email based on the name.
func GetLocalizedEmail(firstName, lastName string) string {
	domains := []string{"example.com", "test.org", "mail.com"}
	// This is not strictly random but deterministic based on name
	domain := domains[(len(firstName)+len(lastName))%len(domains)]
	return fmt.Sprintf("%s.%s@%s", firstName, lastName, domain)
}
