package analyze

import "regexp"

// Rule represents a heuristic for identifying sensitive data.
type Rule struct {
	Name      string
	Regex     *regexp.Regexp
	Generator string
	Provider  string
}

// DefaultRules is the set of heuristic rules used to identify sensitive columns.
var DefaultRules = []Rule{
	{
		Name:      "Email",
		Regex:     regexp.MustCompile(`(?i)email`),
		Generator: "faker",
		Provider:  "email",
	},
	{
		Name:      "First Name",
		Regex:     regexp.MustCompile(`(?i)first_name|fname`),
		Generator: "faker",
		Provider:  "first_name",
	},
	{
		Name:      "Last Name",
		Regex:     regexp.MustCompile(`(?i)last_name|lname`),
		Generator: "faker",
		Provider:  "last_name",
	},
	{
		Name:      "Full Name",
		Regex:     regexp.MustCompile(`(?i)full_name|fullname|name`),
		Generator: "faker",
		Provider:  "full_name",
	},
	{
		Name:      "Phone",
		Regex:     regexp.MustCompile(`(?i)phone|tel|mobile`),
		Generator: "faker",
		Provider:  "phone_number",
	},
	{
		Name:      "Company",
		Regex:     regexp.MustCompile(`(?i)company|organization|org`),
		Generator: "faker",
		Provider:  "company_name",
	},
	{
		Name:      "IP Address",
		Regex:     regexp.MustCompile(`(?i)ip_address|ipv4|ipv6`),
		Generator: "faker",
		Provider:  "ipv4_address",
	},
	{
		Name:      "Address",
		Regex:     regexp.MustCompile(`(?i)address|street|city|zip`),
		Generator: "faker",
		Provider:  "short_text",
	},
}
