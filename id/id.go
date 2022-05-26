// Package id is a simple wrapper around conv. It forces conversion using a human-readable token set where typos should
// occur less frequently: e.g., there are no I's (resembles 1), no O's (resembles 0).
package id

import (
	"strings"

	"github.com/KarelKubat/hrid/conv"
)

const (
	// Tokens are the default runes from which IDs may be constructed.
	Tokens = "0123456789ABCDEFGHKLMNPQRSTUVWXY"
	// StringLen defines the default padding for ID generation.
	StringLen = 14
	// IgnoreCase defines whether ID to number conversions care about casing.
	IgnoreCase = true
	// GroupSize defines the length of groups in generated IDs, for better readability.
	GroupSize = 4
	// Default number of checksum runes to add to a generated ID.
	ChecksumLen = 2
)

// Opts defines the options when constructing an ID converter.
type Opts struct {
	Tokens      string // Tokens to use for conversion: "01" for binary, "0123456789" for decimal, etc.
	StringLen   int    // Minimum length of an ID, which is left-padded with the first token (interpreted as zero).
	IgnoreCase  bool   // When true, casing will be ignored during conversions.
	GroupSize   int    // When non-zero, an ID will be split into space-delimited groups for readability (e.g. "0123 4567").
	ChecksumLen int    // Number of checksum runes to add to an ID, 0 for no checksumming.
}

// ID is the receiver that implements conversions.
type ID struct {
	opts      *Opts
	converter *conv.Conv
}

// New instantiates a converter.
func New(o *Opts) (*ID, error) {
	if o.IgnoreCase {
		o.Tokens = strings.ToUpper(o.Tokens)
	}
	conv, err := conv.New(o.Tokens, o.ChecksumLen)
	if err != nil {
		return nil, err
	}
	return &ID{
		opts:      o,
		converter: conv,
	}, nil
}

// ToRunes converts a uint64 to a slice of runes.
func (id *ID) ToRunes(n uint64) []rune {
	out := id.converter.ToRunes(n)

	// Prepend the first alphabet rune until the desired length is reached.
	for len(out) < id.opts.StringLen {
		out = append([]rune{id.converter.FirstRune()}, out...)
	}

	// Split into groups if requested.
	if id.opts.GroupSize > 0 {
		formatted := []rune{}
		for i := 0; i < len(out); i += id.opts.GroupSize {
			if len(formatted) > 0 {
				formatted = append(formatted, ' ')
			}
			end := i + id.opts.GroupSize
			if end > len(out) {
				end = len(out)
			}
			formatted = append(formatted, out[i:end]...)
		}
		out = formatted
	}
	return out
}

// ToString converts a uint64 to a string.
func (id *ID) ToString(n uint64) string {
	return string(id.ToRunes(n))
}

// ToNr converts a string to a uint64.
func (id *ID) ToNr(s string) (uint64, error) {
	if id.opts.IgnoreCase {
		s = strings.ToUpper(s)
	}
	if id.opts.GroupSize > 0 {
		s = strings.Join(strings.Fields(s), "")
	}
	return id.converter.ToNr(s)
}

var converter *ID

func init() {
	var err error
	converter, err = New(&Opts{
		Tokens:     Tokens,
		StringLen:  StringLen,
		IgnoreCase: IgnoreCase,
		GroupSize:  GroupSize,
	})
	if err != nil {
		panic("failed to construct default converter")
	}
}

// ToString returns the string representation of a uint64, using the defaults.
func ToString(n uint64) string {
	return converter.ToString(n)
}

// ToNr returns the uint64 representation of a string, using the defaults. An error occurs when the string contains
// a rune that is not in the available token set.
func ToNr(s string) (uint64, error) {
	return converter.ToNr(s)
}
