package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/matriphe/postcode"
	"github.com/matriphe/postcode/testdata"
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

func main() {
	var cases []goldenCase
	if err := json.Unmarshal(testdata.GoldenCases, &cases); err != nil {
		fail("decode golden cases: %v", err)
	}

	for _, tt := range cases {
		p, err := postcode.In(tt.Country)
		if err != nil {
			fail("In(%q): %v", tt.Country, err)
		}
		if got := p.IsValid(tt.Input); got != tt.Valid {
			fail("IsValid(%q, %q) = %v, want %v", tt.Country, tt.Input, got, tt.Valid)
		}
		if tt.Outputs == nil {
			continue
		}
		if got := p.Format(tt.Input); got != tt.Outputs.Format {
			fail("Format(%q, %q) = %q, want %q", tt.Country, tt.Input, got, tt.Outputs.Format)
		}
		if got := p.Display(tt.Input); got != tt.Outputs.Display {
			fail("Display(%q, %q) = %q, want %q", tt.Country, tt.Input, got, tt.Outputs.Display)
		}
		if got := p.Code(tt.Input); got != tt.Outputs.Code {
			fail("Code(%q, %q) = %q, want %q", tt.Country, tt.Input, got, tt.Outputs.Code)
		}
		if got := p.FullDisplay(tt.Input); got != tt.Outputs.Full {
			fail("FullDisplay(%q, %q) = %q, want %q", tt.Country, tt.Input, got, tt.Outputs.Full)
		}
		if got := p.Suffix(tt.Input); got != tt.Outputs.Suffix {
			fail("Suffix(%q, %q) = %q, want %q", tt.Country, tt.Input, got, tt.Outputs.Suffix)
		}
	}

	fmt.Printf("checked %d postcode golden cases\n", len(cases))
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
