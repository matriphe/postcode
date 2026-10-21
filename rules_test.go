package postcode

import "testing"

func TestEmbeddedExamplesMatchTheirCountryPatterns(t *testing.T) {
	std := New()
	for countryCode, rule := range std.rules {
		for _, example := range rule.definition.Examples {
			if !rule.pattern.MatchString(example) {
				t.Errorf("%s example %q does not match its pattern", countryCode, example)
			}
		}
	}
}
