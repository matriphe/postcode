// Package postcode validates and formats postal codes by country.
package postcode

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
	"unicode"

	"github.com/tidwall/jsonc"
)

//go:embed data/rules.jsonc
var embeddedRules []byte

var (
	defaultOnce sync.Once
	defaultStd  *Standardizer
	defaultErr  error
)

// Config contains country rules that can be merged with the embedded defaults.
type Config struct {
	Version     int                    `json:"version,omitempty"`
	Attribution string                 `json:"attribution,omitempty"`
	Countries   map[string]CountryRule `json:"countries"`
}

// CountryRule describes the accepted input and its country-specific output forms.
type CountryRule struct {
	Name     string        `json:"name"`
	Pattern  string        `json:"pattern"`
	Sanitize Sanitizer     `json:"sanitize,omitempty"`
	Format   []Replacement `json:"format,omitempty"`
	Primary  []Replacement `json:"primary,omitempty"`
	Display  []Replacement `json:"display,omitempty"`
	Suffix   []Replacement `json:"suffix,omitempty"`
	Code     []Replacement `json:"code,omitempty"`
	Examples []string      `json:"examples,omitempty"`
	Sources  []Source      `json:"sources,omitempty"`
}

// Sanitizer describes characters removed before formatting.
type Sanitizer struct {
	Remove    string `json:"remove,omitempty"`
	Uppercase bool   `json:"uppercase,omitempty"`
}

// Replacement applies a regular-expression replacement to a formatted code.
type Replacement struct {
	Pattern string `json:"pattern"`
	Replace string `json:"replace"`
}

// Source records provenance for a country rule.
type Source struct {
	Name     string `json:"name,omitempty"`
	URL      string `json:"url,omitempty"`
	License  string `json:"license,omitempty"`
	Verified string `json:"verified,omitempty"`
	Note     string `json:"note,omitempty"`
}

// Standardizer selects postcode rules from a private copy of the defaults.
type Standardizer struct {
	rules map[string]*compiledRule
}

type compiledRule struct {
	definition CountryRule
	pattern    *regexp.Regexp
	format     []compiledReplacement
	primary    []compiledReplacement
	display    []compiledReplacement
	suffix     []compiledReplacement
	code       []compiledReplacement
}

type compiledReplacement struct {
	pattern *regexp.Regexp
	replace string
}

// Postcode applies one country's postcode rules.
type Postcode struct {
	rule *compiledRule
}

// New creates an independent standardizer with the embedded country rules.
func New() *Standardizer {
	std := embeddedStandardizer()
	rules := make(map[string]*compiledRule, len(std.rules))
	for code, rule := range std.rules {
		rules[code] = rule
	}
	return &Standardizer{rules: rules}
}

// In returns a handler for the requested country code.
func In(countryCode string) (*Postcode, error) {
	return embeddedStandardizer().In(countryCode)
}

// In returns a handler for the requested country code.
func (s *Standardizer) In(countryCode string) (*Postcode, error) {
	if s == nil {
		return nil, fmt.Errorf("postcode: nil standardizer")
	}
	code := normalizeCountryCode(countryCode)
	rule, ok := s.rules[code]
	if !ok {
		return nil, fmt.Errorf("postcode: unsupported country code %q", countryCode)
	}
	return &Postcode{rule: rule}, nil
}

// WithConfig returns a new standardizer with config rules merged over its rules.
func (s *Standardizer) WithConfig(config Config) (*Standardizer, error) {
	if s == nil {
		return nil, fmt.Errorf("postcode: nil standardizer")
	}
	if config.Version != 0 && config.Version != 1 {
		return nil, fmt.Errorf("postcode: unsupported config version %d", config.Version)
	}
	if len(config.Countries) == 0 {
		return nil, fmt.Errorf("postcode: config has no countries")
	}

	rules := make(map[string]*compiledRule, len(s.rules)+len(config.Countries))
	for code, rule := range s.rules {
		rules[code] = rule
	}
	for countryCode, definition := range config.Countries {
		code := normalizeCountryCode(countryCode)
		if !validCountryCode(code) {
			return nil, fmt.Errorf("postcode: invalid country code %q", countryCode)
		}
		compiled, err := compileRule(definition)
		if err != nil {
			return nil, fmt.Errorf("postcode: country %s: %w", code, err)
		}
		rules[code] = compiled
	}
	return &Standardizer{rules: rules}, nil
}

// LoadConfig decodes a JSONC configuration.
func LoadConfig(reader io.Reader) (Config, error) {
	var config Config
	if reader == nil {
		return config, fmt.Errorf("postcode: nil config reader")
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return config, fmt.Errorf("postcode: read config: %w", err)
	}
	clean := jsonc.ToJSON(data)
	decoder := json.NewDecoder(bytes.NewReader(clean))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&config); err != nil {
		return config, fmt.Errorf("postcode: decode config: %w", err)
	}
	if config.Version != 0 && config.Version != 1 {
		return config, fmt.Errorf("postcode: unsupported config version %d", config.Version)
	}
	if len(config.Countries) == 0 {
		return config, fmt.Errorf("postcode: config has no countries")
	}
	return config, nil
}

// IsValid reports whether input matches the country's postcode format.
func (p *Postcode) IsValid(input string) bool {
	return p != nil && p.rule != nil && input != "" && p.rule.pattern.MatchString(input)
}

// Format returns the primary postcode in compact form, or an empty string if invalid.
func (p *Postcode) Format(input string) string {
	formatted := p.canonical(input)
	if formatted == "" {
		return ""
	}
	return applyReplacements(formatted, p.rule.primary)
}

// Display returns the country's display form, or an empty string if invalid.
func (p *Postcode) Display(input string) string {
	formatted := p.Format(input)
	if formatted == "" {
		return ""
	}
	return applyReplacements(formatted, p.rule.display)
}

// FullDisplay returns the complete postcode in its conventional display form.
// For countries without an extended postcode, this is the same as Display.
func (p *Postcode) FullDisplay(input string) string {
	formatted := p.canonical(input)
	if formatted == "" {
		return ""
	}
	return applyReplacements(formatted, p.rule.display)
}

// Suffix returns an optional country-specific postcode suffix, or an empty
// string when the postcode has no suffix or the input is invalid.
func (p *Postcode) Suffix(input string) string {
	formatted := p.canonical(input)
	if formatted == "" || len(p.rule.suffix) == 0 {
		return ""
	}
	return applyReplacements(formatted, p.rule.suffix)
}

// Code returns the canonical postcode without any country-specific prefix.
func (p *Postcode) Code(input string) string {
	formatted := p.canonical(input)
	if formatted == "" {
		return ""
	}
	return applyReplacements(formatted, p.rule.code)
}

func (p *Postcode) canonical(input string) string {
	if p == nil || p.rule == nil || input == "" {
		return ""
	}
	clean := sanitize(input, p.rule.definition.Sanitize)
	if clean == "" || (!p.rule.pattern.MatchString(input) && !p.rule.pattern.MatchString(clean)) {
		return ""
	}
	return applyReplacements(clean, p.rule.format)
}

func embeddedStandardizer() *Standardizer {
	defaultOnce.Do(func() {
		data, err := io.ReadAll(bytes.NewReader(embeddedRules))
		if err != nil {
			defaultErr = err
			return
		}
		clean := jsonc.ToJSON(data)
		var config Config
		if err := json.Unmarshal(clean, &config); err != nil {
			defaultErr = err
			return
		}
		defaultStd, defaultErr = compileConfig(config)
	})
	if defaultErr != nil {
		panic("postcode: invalid embedded rules: " + defaultErr.Error())
	}
	return defaultStd
}

func compileConfig(config Config) (*Standardizer, error) {
	if len(config.Countries) == 0 {
		return nil, fmt.Errorf("config has no countries")
	}
	rules := make(map[string]*compiledRule, len(config.Countries))
	for countryCode, definition := range config.Countries {
		code := normalizeCountryCode(countryCode)
		if !validCountryCode(code) {
			return nil, fmt.Errorf("invalid country code %q", countryCode)
		}
		compiled, err := compileRule(definition)
		if err != nil {
			return nil, fmt.Errorf("country %s: %w", code, err)
		}
		rules[code] = compiled
	}
	return &Standardizer{rules: rules}, nil
}

func compileRule(definition CountryRule) (*compiledRule, error) {
	if strings.TrimSpace(definition.Name) == "" {
		return nil, fmt.Errorf("name is required")
	}
	if strings.TrimSpace(definition.Pattern) == "" {
		return nil, fmt.Errorf("pattern is required")
	}
	pattern, err := regexp.Compile("(?i)" + definition.Pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}
	format, err := compileReplacements(definition.Format)
	if err != nil {
		return nil, fmt.Errorf("format replacement: %w", err)
	}
	primary, err := compileReplacements(definition.Primary)
	if err != nil {
		return nil, fmt.Errorf("primary replacement: %w", err)
	}
	display, err := compileReplacements(definition.Display)
	if err != nil {
		return nil, fmt.Errorf("display replacement: %w", err)
	}
	suffix, err := compileReplacements(definition.Suffix)
	if err != nil {
		return nil, fmt.Errorf("suffix replacement: %w", err)
	}
	code, err := compileReplacements(definition.Code)
	if err != nil {
		return nil, fmt.Errorf("code replacement: %w", err)
	}
	return &compiledRule{definition: definition, pattern: pattern, format: format, primary: primary, display: display, suffix: suffix, code: code}, nil
}

func compileReplacements(replacements []Replacement) ([]compiledReplacement, error) {
	compiled := make([]compiledReplacement, 0, len(replacements))
	for _, replacement := range replacements {
		pattern, err := regexp.Compile(replacement.Pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern %q: %w", replacement.Pattern, err)
		}
		compiled = append(compiled, compiledReplacement{pattern: pattern, replace: replacement.Replace})
	}
	return compiled, nil
}

func applyReplacements(input string, replacements []compiledReplacement) string {
	for _, replacement := range replacements {
		input = replacement.pattern.ReplaceAllString(input, replacement.replace)
	}
	return input
}

func sanitize(input string, rule Sanitizer) string {
	clean := strings.Map(func(r rune) rune {
		if strings.ContainsRune(rule.Remove, r) || (unicode.IsSpace(r) && strings.ContainsRune(rule.Remove, ' ')) {
			return -1
		}
		return r
	}, input)
	if rule.Uppercase {
		clean = strings.ToUpper(clean)
	}
	return clean
}

func normalizeCountryCode(countryCode string) string {
	return strings.ToUpper(strings.TrimSpace(countryCode))
}

func validCountryCode(countryCode string) bool {
	if len(countryCode) != 2 {
		return false
	}
	for _, r := range countryCode {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}
