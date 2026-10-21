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
}
```

`IsValid` checks whether a value matches the country's postcode format. It does not check whether the postcode is assigned or deliverable. `Format`, `Display`, and `Code` return an empty string for invalid input. `Format` returns the compact canonical form, `Display` returns the conventional display form, and `Code` removes a country-specific prefix where configured. For example, Luxembourg's `L-1234` becomes `L1234` with `Format`, `L-1234` with `Display`, and `1234` with `Code`.

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
