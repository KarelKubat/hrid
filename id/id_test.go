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

func TestFormattingWithoutChecksum(t *testing.T) {
	for _, test := range []struct {
		groupSize  int
		wantString string
	}{
		{
			groupSize:  0,
			wantString: "123456789",
		},
		{
			groupSize:  1,
			wantString: "1 2 3 4 5 6 7 8 9",
		},
		{
			groupSize:  2,
			wantString: "12 34 56 78 9",
		},
		{
			groupSize:  3,
			wantString: "123 456 789",
		},
		{
			groupSize:  4,
			wantString: "1234 5678 9",
		},
		{
			groupSize:  5,
			wantString: "12345 6789",
		},
		{
			groupSize:  6,
			wantString: "123456 789",
		},
		{
			groupSize:  7,
			wantString: "1234567 89",
		},
		{
			groupSize:  8,
			wantString: "12345678 9",
		},
		{
			groupSize:  9,
			wantString: "123456789",
		},
		{
			groupSize:  666,
			wantString: "123456789",
		},
	} {
		id, err := New(&Opts{
			Alphabet:    "0123456789",
			StringLen:   0,
			IgnoreCase:  false,
			GroupSize:   test.groupSize,
			ChecksumLen: 0,
		})
		if err != nil {
			t.Fatalf("New() = _,%v, need nil error", err)
		}
		if gotString := id.ToString(uint64(123456789)); gotString != test.wantString {
			t.Errorf("id.Tostring() using grupsize %v = %v, want %v", test.groupSize, gotString, test.wantString)
		}
	}
}
