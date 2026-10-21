# postcode

`postcode` validates, formats, and displays postal codes using country-specific rules embedded in the package. It works offline and supports Go 1.16 and later.

```sh
go get github.com/matriphe/postcode
```

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

Country codes are case-insensitive ISO alpha-2 codes. `IsValid` checks format only; it does not verify that a postcode is assigned or deliverable. Formatting methods return an empty string for invalid input.

`Format` returns compact canonical form, `Display` returns the conventional display form, and `Code` removes a country-specific prefix where configured. For example, Luxembourg's `L-1234` becomes `L1234` with `Format`, `L-1234` with `Display`, and `1234` with `Code`.

Custom JSONC rules can be loaded and merged into an independent standardizer:

```go
import (
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
	_, err = std.In("ZZ")
	return err
}
```

The embedded JSONC rules currently cover 179 countries and territories with available postcode patterns. Each rule records its source and verification date. A missing rule means the package cannot validate that country's format.

The package source is licensed under MIT. The rule data includes material derived from Google `libaddressinput`, licensed CC BY 4.0; see [data/README.md](data/README.md) and the source records in [data/rules.jsonc](data/rules.jsonc).
