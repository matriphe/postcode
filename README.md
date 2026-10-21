# postcode

Validate, format, and display postal codes with country-specific rules embedded in the package. Works offline and supports Go 1.16 and later.

## Installation

```sh
go get github.com/matriphe/postcode
```

## Usage

Select a country with its case-insensitive ISO alpha-2 code, then use its postcode handler:

```go
package main

import (
	"fmt"
	"log"

	"github.com/matriphe/postcode"
)

func main() {
	p, err := postcode.In("DE")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(p.IsValid("10557"))       // true
	fmt.Println(p.Format(" 1 0 5 5 7 ")) // 10557
	fmt.Println(p.Display("10557"))       // 10557
	fmt.Println(p.Code("10557"))          // 10557
}
```

`IsValid` checks whether a value matches the country's postcode format. It does not check whether the postcode is assigned or deliverable. Formatting methods return an empty string for invalid input. `Format` returns the primary postcode in compact form, `Display` returns its conventional display form, and `Code` returns the complete canonical code without country-specific prefixes. Postcodes are strings, so leading zeroes are preserved.

For the US, `IsValid` accepts a five-digit ZIP Code or ZIP+4 with a hyphen, space, or no separator. `Format` and `Display` return the five-digit primary ZIP Code. `FullDisplay` returns the complete ZIP+4 form when present, and `Suffix` returns its final four digits. For other countries, `FullDisplay` matches `Display` and `Suffix` is empty. For example, Luxembourg's `L-1234` becomes `L1234` with `Format`, `L-1234` with `Display` and `FullDisplay`, and `1234` with `Code`.

Country-specific examples:

```go
examples := []struct {
	country string
	input   string
}{
	{"LU", "1234"},
	{"GB", "ec1y 8sy"},
	{"NL", "1234 AB"},
	{"US", "00501"},
	{"US", "00501-1234"},
	{"US", "005011234"},
}

for _, example := range examples {
	p, err := postcode.In(example.country)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s: valid=%t format=%q display=%q full=%q suffix=%q code=%q\n",
		example.country,
		p.IsValid(example.input),
		p.Format(example.input),
		p.Display(example.input),
		p.FullDisplay(example.input),
		p.Suffix(example.input),
		p.Code(example.input),
	)
}
```

This prints:

```text
LU: valid=true format="L1234" display="L-1234" full="L-1234" suffix="" code="1234"
GB: valid=true format="EC1Y8SY" display="EC1Y 8SY" full="EC1Y 8SY" suffix="" code="EC1Y8SY"
NL: valid=true format="1234AB" display="1234 AB" full="1234 AB" suffix="" code="1234AB"
US: valid=true format="00501" display="00501" full="00501" suffix="" code="00501"
US: valid=true format="00501" display="00501" full="00501-1234" suffix="1234" code="005011234"
US: valid=true format="00501" display="00501" full="00501-1234" suffix="1234" code="005011234"
```

The country-specific snippet uses the `fmt` and `log` imports from the quick-start example.

Custom JSONC rules can be merged into an independent standardizer created with `postcode.New()`:

```go
package main

import (
	"fmt"
	"strings"

	"github.com/matriphe/postcode"
)

func custom() error {
	config, err := postcode.LoadConfig(strings.NewReader(`{
  "countries": {
    "ZZ": {"name": "Example", "pattern": "^[0-9]{4}$"}
  }
}`))
	if err != nil {
		return err
	}

	std, err := postcode.New().WithConfig(config)
	if err != nil {
		return err
	}
	p, err := std.In("ZZ")
	if err != nil {
		return err
	}
	fmt.Println(p.IsValid("1234")) // true
	return nil
}
```

A custom config can add or override country rules without changing the package's default rules.

## Data/Rules

Country postcode patterns and examples are embedded in the package, so validation works offline and does not require a separate data download. The rules currently cover 179 countries and territories for which a usable pattern was available. Each rule records its source and verification date; countries without a usable pattern are omitted rather than assigned a guessed rule.

The data includes material derived from Google `libaddressinput` and country-specific overrides. See [data/README.md](data/README.md) for attribution and coverage details, and [data/rules.jsonc](data/rules.jsonc) for the rules and their sources.

## License

The Go package source code is licensed under the [MIT License](LICENSE). The embedded rule data includes material derived from Google `libaddressinput`, licensed under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/). The data's attribution and source information are documented in [data/README.md](data/README.md). Distributors of applications that include the embedded rules should preserve that attribution.
