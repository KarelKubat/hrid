// Package er represents the errors that hrid may throw.
package er

import (
	"fmt"
)

// Code designates the type of an error.
type Code int

const (
	None Code = iota // Keep at first slot with value zero

	AlphabetTooShortError
	TokenRepeatsError
	IDTooShortError
	ChecksumError
	NoSuchTokenError

	ZZLastUnused // Keep at last slot for test coverage
)

// String stringifies a Code.
func (c Code) String() string {
	return []string{
		"None",

		"AlphabetTooShortError",
		"TokenRepeatsError",
		"IDTooShortError",
		"ChecksumError",
		"NoSuchTokenError",
	}[c]
}

// Err contains the error code and its description.
type Err struct {
	Code Code
	Msg  string
}

// New returns an Err given a code and a description.
func New(c Code, msg string) *Err {
	return &Err{
		Code: c,
		Msg:  msg,
	}
}

// Newf is like New but allows printf-like expansion.
func Newf(c Code, format string, args ...interface{}) *Err {
	return New(c, fmt.Sprintf(format, args...))
}

// Error satisfies the error interface.
func (e *Err) Error() string {
	return fmt.Sprintf("%v: %v", e.Code, e.Msg)
}
