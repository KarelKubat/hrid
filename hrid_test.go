package main

import "testing"

func TestMain(t *testing.T) {
	// Make sure that main's worker works.
	// This is similar to the CLI `hrid 12`, which should output something like `00000 00000 00CCR` and we don't care
	// about the actual output, as long as the main binary works we're fine here. Other tests check the conversions.
	hrid([]string{"12"})
}
