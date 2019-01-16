package server

import (
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const (
	maxSearchLen = 64
)

// https://rosettacode.org/wiki/Strip_control_codes_and_extended_characters_from_a_string#Go
func ParseSearchText(text string) string {
	isOk := func(r rune) bool {
		return r < 32 || r >= 127
	}
	// The isOk filter is such that there is no need to chain to norm.NFC
	t := transform.Chain(norm.NFKD, transform.RemoveFunc(isOk))
	// This Transformer could also trivially be applied as an io.Reader
	// or io.Writer filter to automatically do such filtering when reading
	// or writing data anywhere.
	text, _, _ = transform.String(t, text)
	if len(text) > maxSearchLen {
		text = text[:maxSearchLen]
	}
	return text
}
