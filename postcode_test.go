package postcode_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/matriphe/postcode"
)

type goldenCase struct {
	Country string        `json:"country"`
	Input   string        `json:"input"`
	Valid   bool          `json:"valid"`
	Outputs *goldenOutput `json:"outputs,omitempty"`
}

type goldenOutput struct {
	Format  string `json:"format"`
	Display string `json:"display"`
	Code    string `json:"code"`
	Full    string `json:"full"`
	Suffix  string `json:"suffix"`
}

func loadGoldenCases(t *testing.T) []goldenCase {
	t.Helper()
	data, err := os.ReadFile("ci/testdata/postcode-golden.json")
	if err != nil {
		t.Fatalf("read golden cases: %v", err)
	}
	var cases []goldenCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("decode golden cases: %v", err)
	}
	return cases
}

func TestGoldenCases(t *testing.T) {
	for _, tt := range loadGoldenCases(t) {
		t.Run(tt.Country+"/"+tt.Input, func(t *testing.T) {
			p, err := postcode.In(tt.Country)
			if err != nil {
				t.Fatalf("In(%q): %v", tt.Country, err)
			}
			if got := p.IsValid(tt.Input); got != tt.Valid {
				t.Errorf("IsValid(%q) = %v, want %v", tt.Input, got, tt.Valid)
			}
			if tt.Outputs != nil {
				if got := p.Format(tt.Input); got != tt.Outputs.Format {
					t.Errorf("Format(%q) = %q, want %q", tt.Input, got, tt.Outputs.Format)
				}
				if got := p.Display(tt.Input); got != tt.Outputs.Display {
					t.Errorf("Display(%q) = %q, want %q", tt.Input, got, tt.Outputs.Display)
				}
				if got := p.Code(tt.Input); got != tt.Outputs.Code {
					t.Errorf("Code(%q) = %q, want %q", tt.Input, got, tt.Outputs.Code)
				}
				if got := p.FullDisplay(tt.Input); got != tt.Outputs.Full {
					t.Errorf("FullDisplay(%q) = %q, want %q", tt.Input, got, tt.Outputs.Full)
				}
				if got := p.Suffix(tt.Input); got != tt.Outputs.Suffix {
					t.Errorf("Suffix(%q) = %q, want %q", tt.Input, got, tt.Outputs.Suffix)
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
	for country, input := range map[string]string{"DE": "1055", "US": "00501-123"} {
		p, err := postcode.In(country)
		if err != nil {
			t.Fatal(err)
		}
		for name, got := range map[string]string{
			"Format":      p.Format(input),
			"Display":     p.Display(input),
			"Code":        p.Code(input),
			"FullDisplay": p.FullDisplay(input),
			"Suffix":      p.Suffix(input),
		} {
			if got != "" {
				t.Errorf("%s(%q, %s) = %q, want empty", name, input, country, got)
			}
		}
	}
}
