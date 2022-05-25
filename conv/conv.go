// Package conv allows one to define the available tokens for numeric conversions, and then to convert
// between strings and numbers.
package conv

import (
	"errors"
	"fmt"
	"strings"
)

// Set is the receiver that implements ToString and Touint64.
type Conv struct {
	tokens   string
	tokenLen int
}

// New returns a new Conv. The input is e.g. for decimal conversions: "0123456789", for binary: "01", etc.
func New(tokens string) (*Conv, error) {
	if tokens == "" {
		return nil, errors.New("conversion tokens may not be an empty string")
	}
	partsMap := map[string]struct{}{}
	for _, part := range strings.Split(tokens, "") {
		if _, ok := partsMap[part]; ok {
			return nil, fmt.Errorf("%v repeats in tokens %q", part, tokens)
		}
		partsMap[part] = struct{}{}
	}

	return &Conv{
		tokens:   tokens,
		tokenLen: len(tokens),
	}, nil
}

// ToString converts a uint64 to its string representation.
func (a *Conv) ToString(nr uint64) string {
	reversed := ""
	for nr > 0 {
		remainder := nr % uint64(a.tokenLen)
		nr = nr / uint64(a.tokenLen)
		reversed += string([]byte{a.tokens[remainder]})
	}
	if reversed == "" {
		reversed = string([]byte{a.tokens[0]})
	}

	runes := make([]rune, len(reversed))
	n := 0
	for _, r := range reversed {
		runes[n] = r
		n++
	}
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
	}
	return string(runes)
}

// Touint64 converts a string to its numeric representation. An error occurs when the string contains runes that
// are not in the available tokens.
func (a *Conv) ToNr(s string) (uint64, error) {
	out := uint64(0)
	pwr := 0
	for i := len(s) - 1; i >= 0; i-- {
		token := s[i : i+1]
		index := strings.Index(a.tokens, token)
		if index == -1 {
			return 0, fmt.Errorf("token %q not in %q", token, a.tokens)
		}
		// Can't use math.Pow() because of the float64 conversions. The below fails at large uint64 values.
		// out += uint64(math.Pow(float64(a.tokenLen), float64(pwr)) * float64(index))
		out += intPow(a.tokenLen, pwr) * uint64(index)
		pwr += 1
	}
	return out, nil
}

func intPow(m, e int) uint64 {
	if e == 0 {
		return 1
	}
	out := uint64(m)
	for i := 2; i <= e; i++ {
		out *= uint64(m)
	}
	return out
}
