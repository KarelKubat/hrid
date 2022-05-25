// Package id is a simple wrapper around conv. It forces conversion using a human-readable token set where typos should
// occur less frequently: e.g., there are no I's (resembles 1), no O's (resembles 0).
package id

import (
	"strings"

	"github.com/KarelKubat/hrid/conv"
)

const (
	// Tokens are the default runes from which IDs may be constructed.
	Tokens = "0123456789ABCDEFGHKLMNPQRSTUVWXYZ"
	// StringLen defines the default padding for ID generation.
	StringLen = 20
	// IgnoreCase defines whether ID to number conversions care about casing.
	IgnoreCase = true
	// GroupSize defines the length of groups in generated IDs, for better readability.
	GroupSize = 4
)

// Opts defines the options when constructing an ID converter.
type Opts struct {
	Tokens     string // Tokens to use for conversion: "01" for binary, "0123456789" for decimal, etc.
	StringLen  int    // Length of an ID, which is left-padded with the first token (interpreted as zero).
	IgnoreCase bool   // When true, casing will be ignored during conversions.
	GroupSize  int    // When non-zero, an ID will be split into space-delimited groups for readability (e.g. "0123 4567").
}

type ID struct {
	opts      *Opts
	converter *conv.Conv
}

func New(o *Opts) (*ID, error) {
	if o.IgnoreCase {
		o.Tokens = strings.ToUpper(o.Tokens)
	}
	conv, err := conv.New(o.Tokens)
	if err != nil {
		return nil, err
	}
	return &ID{
		opts:      o,
		converter: conv,
	}, nil
}

func (id *ID) ToString(n uint64) string {
	out := id.converter.ToString(n)
	for len(out) < id.opts.StringLen {
		out = id.opts.Tokens[0:1] + out
	}
	if id.opts.GroupSize > 0 {
		formatted := ""
		for i := 0; i < len(out); i += id.opts.GroupSize {
			end := i + id.opts.GroupSize
			if end > len(out) {
				end = len(out)
			}
			formatted += out[i:end] + " "
		}
		out = formatted
	}
	return out
}

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

// ToString returns the string representation of a uint64, given the available tokens.
func ToString(n uint64) string {
	return converter.ToString(n)
}

// ToNr returns the uint64 representation of a string, given the available tokens. An error occurs when the string contains
// a rune that is not in the available token set.
func ToNr(s string) (uint64, error) {
	return converter.ToNr(s)
}
