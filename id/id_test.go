package id

import "testing"

func TestConversions(t *testing.T) {
	for _, n := range []uint64{
		0, 1, 2, 3, 4, 5,
		42,
		1234567890, 987654321,
	} {
		s := ToString(n)
		gotNr, err := ToNr(s)
		if err != nil {
			t.Fatalf("ToNr(%q) = _,%q, need nil error", s, err)
		}
		if gotNr != n {
			t.Errorf("ToID(%v) = %q, but ToNr(%q) = %v, mismatch", n, s, s, gotNr)
		}
	}
}

func TestError(t *testing.T) {
	id := "012 . 345"
	_, err := ToNr(id)
	if err == nil {
		t.Fatalf("ToNr(%q) returns nil error", id)
	}
}
