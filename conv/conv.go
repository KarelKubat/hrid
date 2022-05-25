// Package conv allows one to define the available tokens for numeric conversions, and then to convert
// between strings and numbers.
package conv

import (
	"errors"
	"fmt"
)

// Set is the receiver that implements ToString and Touint64.
type Conv struct {
	tokens     []rune
	tokenIndex map[rune]int
	tokenLen   int
}

// New returns a new Conv. The input is e.g. for decimal conversions: "0123456789", for binary: "01", etc.
func New(alphabet string) (*Conv, error) {
	if alphabet == "" {
		return nil, errors.New("conversion tokens may not be an empty string")
	}
	tokens := []rune(alphabet)
	tokenIndex := map[rune]int{}
	for i, part := range tokens {
		if _, ok := tokenIndex[part]; ok {
			return nil, fmt.Errorf("%v repeats in tokens %q", part, tokens)
		}
		tokenIndex[part] = i
	}

	return &Conv{
		tokens:     tokens,
		tokenIndex: tokenIndex,
		tokenLen:   len(tokens),
	}, nil
}

// FirstRune returns the first rune of the tokens alphabet.
func (a *Conv) FirstRune() rune {
	return a.tokens[0]
}

// ToRunes converts a uint64 to runes representation.
func (a *Conv) ToRunes(nr uint64) []rune {
	reversed := []rune{}
	for nr > 0 {
		remainder := nr % uint64(a.tokenLen)
		nr = nr / uint64(a.tokenLen)
		reversed = append(reversed, a.tokens[remainder])
	}
	if len(reversed) == 0 {
		reversed = []rune{a.tokens[0]}
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
	return runes
}

// ToString converts a uint64 to a string representation.
func (a *Conv) ToString(nr uint64) string {
	return string(a.ToRunes(nr))
}

// ToNr converts a string to its numeric representation. An error occurs when the string contains runes that
// are not in the available tokens.
func (a *Conv) ToNr(s string) (uint64, error) {
	out := uint64(0)
	pwr := 0
	tokens := []rune(s)
	for i := len(tokens) - 1; i >= 0; i-- {
		token := tokens[i]
		index, ok := a.tokenIndex[token]
		if !ok {
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
