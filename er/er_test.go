package er

import "testing"

func TestString(t *testing.T) {
	// Test coverage that all codes between None and ZZLastUnused are stringable.
	for c := None; c < ZZLastUnused; c++ {
		_ = c.String()
	}
}
