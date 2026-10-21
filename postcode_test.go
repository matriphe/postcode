package postcode_test

import (
	"strings"
	"testing"

	"github.com/matriphe/postcode"
)

func TestCountrySpecificBehavior(t *testing.T) {
	tests := []struct {
		country string
		input   string
		valid   bool
		format  string
		display string
		code    string
	}{
		{country: "de", input: "10557", valid: true, format: "10557", display: "10557", code: "10557"},
		{country: "DE", input: "1055", valid: false},
		{country: "DE", input: "1 0 5 5 7", valid: false, format: "10557", display: "10557", code: "10557"},
		{country: "NO", input: "1055", valid: true, format: "1055", display: "1055", code: "1055"},
		{country: "NO", input: "105", valid: true, format: "0105", display: "0105", code: "0105"},
		{country: "NL", input: "1234 AB", valid: true, format: "1234AB", display: "1234 AB", code: "1234AB"},
		{country: "NL", input: "1234AB", valid: true, format: "1234AB", display: "1234 AB", code: "1234AB"},
		{country: "NL", input: "12 34 AB", valid: false},
		{country: "GB", input: "ec1y 8sy", valid: true, format: "EC1Y8SY", display: "EC1Y 8SY", code: "EC1Y8SY"},
		{country: "GB", input: "GIR 0AA", valid: true, format: "GIR0AA", display: "GIR 0AA", code: "GIR0AA"},
		{country: "GB", input: "ZZ1 1ZZ", valid: false},
		{country: "LU", input: "1234", valid: true, format: "L1234", display: "L-1234", code: "1234"},
		{country: "LU", input: "L1234", valid: true, format: "L1234", display: "L-1234", code: "1234"},
		{country: "LU", input: "L-1234", valid: true, format: "L1234", display: "L-1234", code: "1234"},
		{country: "LU", input: "L-123", valid: false},
		{country: "CA", input: "K1A 0B1", valid: true, format: "K1A0B1", display: "K1A 0B1", code: "K1A0B1"},
		{country: "IE", input: "D02 X285", valid: true, format: "D02X285", display: "D02 X285", code: "D02X285"},
		{country: "US", input: "00501-1234", valid: true, format: "005011234", display: "00501-1234", code: "005011234"},
	}

	for _, tt := range tests {
		t.Run(tt.country+"/"+tt.input, func(t *testing.T) {
			p, err := postcode.In(tt.country)
			if err != nil {
				t.Fatalf("In(%q): %v", tt.country, err)
			}
			if got := p.IsValid(tt.input); got != tt.valid {
				t.Errorf("IsValid(%q) = %v, want %v", tt.input, got, tt.valid)
			}
			if tt.format != "" {
				if got := p.Format(tt.input); got != tt.format {
					t.Errorf("Format(%q) = %q, want %q", tt.input, got, tt.format)
				}
				if got := p.Display(tt.input); got != tt.display {
					t.Errorf("Display(%q) = %q, want %q", tt.input, got, tt.display)
				}
				if got := p.Code(tt.input); got != tt.code {
					t.Errorf("Code(%q) = %q, want %q", tt.input, got, tt.code)
				}
			}
		})
	}
}

func TestUnknownCountry(t *testing.T) {
	if _, err := postcode.In("XX"); err == nil {
		t.Fatal("In(\"XX\") returned no error")
	}
}

func TestLoadJSONCAndMergeWithoutChangingDefaults(t *testing.T) {
	config, err := postcode.LoadConfig(strings.NewReader(`{
		// Overrides the default German format for this standardizer only.
		"countries": {
			"DE": {
				"name": "Custom Germany",
				"pattern": "^[A-Z]{2}[0-9]{3}$",
				"sanitize": {"uppercase": true}
			}
		}
	}`))
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}

	custom, err := postcode.New().WithConfig(config)
	if err != nil {
		t.Fatalf("WithConfig: %v", err)
	}
	customDE, err := custom.In("de")
	if err != nil {
		t.Fatalf("custom.In: %v", err)
	}
	if !customDE.IsValid("AB123") || customDE.IsValid("10557") {
		t.Fatal("custom rule did not override DE in the independent standardizer")
	}

	defaultDE, err := postcode.In("DE")
	if err != nil {
		t.Fatalf("postcode.In: %v", err)
	}
	if !defaultDE.IsValid("10557") || defaultDE.IsValid("AB123") {
		t.Fatal("custom configuration changed the embedded defaults")
	}
}

func TestInvalidInputProducesEmptyStrings(t *testing.T) {
	p, err := postcode.In("DE")
	if err != nil {
		t.Fatal(err)
	}
	for name, got := range map[string]string{
		"Format":  p.Format("1055"),
		"Display": p.Display("1055"),
		"Code":    p.Code("1055"),
	} {
		if got != "" {
			t.Errorf("%s(invalid) = %q, want empty", name, got)
		}
	}
}
