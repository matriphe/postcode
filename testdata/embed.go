// Package testdata embeds the shared postcode cases for the CI smoke app.
package testdata

import _ "embed"

// GoldenCases contains the cases shared by unit and container tests.
//
//go:embed postcode-golden.json
var GoldenCases []byte
