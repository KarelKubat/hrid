// Package conv allows one to define the available tokens for numeric conversions, and then to convert
// between strings and numbers.
package conv

import (
	"errors"
	"fmt"
)

// Set is the receiver that implements ToString and Touint64.
type Conv struct {
	tokens      []rune
	checksumLen int
	tokenIndex  map[rune]int
	tokenLen    int
}

// New returns a new Conv. The input is e.g. for decimal conversions: "0123456789", for binary: "01", etc.
func New(alphabet string, checksumLen int) (*Conv, error) {
	if alphabet == "" {
		return nil, errors.New("conv.New: conversion tokens may not be an empty string")
	}
	tokens := []rune(alphabet)
	tokenIndex := map[rune]int{}
	for i, part := range tokens {
		if _, ok := tokenIndex[part]; ok {
			return nil, fmt.Errorf("conv.New: %v repeats in tokens %q", string(part), string(tokens))
		}
		tokenIndex[part] = i
	}

	return &Conv{
		tokens:      tokens,
		checksumLen: checksumLen,
		tokenIndex:  tokenIndex,
		tokenLen:    len(tokens),
	}, nil
}

// FirstRune returns the first rune of the tokens alphabet.
func (a *Conv) FirstRune() rune {
	return a.tokens[0]
}

// ToRunes converts a uint64 to runes representation and adds checksum runes if so requested.
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

	// Add checksumming of so requested.
	for i := 0; i < a.checksumLen; i++ {
		runes = append(runes, a.checksum(runes))
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
	tokens := []rune(s)

	// Verify checksumming.
	if len(tokens) <= a.checksumLen {
		return 0, fmt.Errorf("conv.ToNr: %q doesn't accomodate %v checksum runes", s, a.checksumLen)
	}
	for i := 0; i < a.checksumLen; i++ {
		gotCs := tokens[len(tokens)-1:][0]
		tokens = tokens[:len(tokens)-1]
		wantCs := a.checksum(tokens)
		if gotCs != wantCs {
			return 0, fmt.Errorf("conv.ToNr: checksum error at %v, expected %v", string(gotCs), string(wantCs))
		}
	}

	// Convert to a number.
	out := uint64(0)
	pwr := 0
	for i := len(tokens) - 1; i >= 0; i-- {
		token := tokens[i]
		index, ok := a.tokenIndex[token]
		if !ok {
			return 0, fmt.Errorf("conv.ToNr: token %v not in %q", string(token), string(a.tokens))
		}
		// Can't use math.Pow() because of the float64 conversions. The below fails at large uint64 values.
		// out += uint64(math.Pow(float64(a.tokenLen), float64(pwr)) * float64(index))
		out += intPow(a.tokenLen, pwr) * uint64(index)
		pwr += 1
	}
	return out, nil
}

// intPow is a helper to compute m to the power of e.
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

// checksum is a helper to compute the checksum of a slice of runes.
func (a *Conv) checksum(runes []rune) rune {
	cs := 0
	for _, r := range runes {
		cs += a.tokenIndex[r]
		cs %= len(a.tokens)
	}
	return a.tokens[cs]
}
