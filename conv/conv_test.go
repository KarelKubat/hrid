package conv

import (
	"testing"
)

func TestIntPow(t *testing.T) {
	for _, test := range []struct {
		mantissa, exponent int
		wantResult         uint64
	}{
		{12, 0, 1},
		{12, 1, 12},
		{12, 2, 144},

		{2, 2, 4},
		{2, 3, 8},
		{2, 4, 16},
		{2, 5, 32},
	} {
		gotResult := intPow(test.mantissa, test.exponent)
		if gotResult != test.wantResult {
			t.Errorf("intPow(%v,%v) = %v but want %v", test.mantissa, test.exponent, gotResult, test.wantResult)
		}
	}
}

func TestToNumber(t *testing.T) {
	a, err := New("0123456789")
	if err != nil {
		t.Fatalf("New(0-9) returned unexpected error %v", err)
	}
	for _, test := range []struct {
		s       string
		wantNr  uint64
		wantErr bool
	}{
		{
			s:      "0",
			wantNr: uint64(0),
		},
		{
			s:      "00000",
			wantNr: uint64(0),
		},
		{
			s:      "12",
			wantNr: uint64(12),
		},
		{
			s:      "123456789",
			wantNr: uint64(123456789),
		},
		{
			s:      "0987654321",
			wantNr: uint64(987654321),
		},
		{
			s:       "09876 x 54321",
			wantErr: true,
		},
	} {
		gotNr, gotErr := a.ToNr(test.s)
		switch {
		case gotErr == nil && test.wantErr:
			t.Errorf("a.ToNr(%q) = _,nil, want error", test.s)
		case gotErr != nil && !test.wantErr:
			t.Errorf("a.ToNr(%q) = _,%q, want no error", test.s, gotErr)
		case gotErr == nil && !test.wantErr && gotNr != test.wantNr:
			t.Errorf("a.ToNr(%q) = %v,_, want nr %v", test.s, gotNr, test.wantNr)
		}
	}

	a, err = New("01")
	if err != nil {
		t.Fatalf("New(01) returned unexpected error %v", err)
	}
	for _, test := range []struct {
		s      string
		wantNr uint64
	}{
		{
			s:      "0",
			wantNr: 0,
		},
		{
			s:      "1",
			wantNr: 1,
		},
		{
			s:      "10",
			wantNr: 2,
		},
		{
			s:      "11",
			wantNr: 3,
		},
		{
			s:      "100",
			wantNr: 4,
		},
	} {
		gotNr, err := a.ToNr(test.s)
		if err != nil {
			t.Fatalf("a.ToNr(%q) = _,%q, need nil error", test.s, err)
		}
		if gotNr != test.wantNr {
			t.Errorf("a.ToNr(%q) = %v,_, want nr %v", test.s, gotNr, test.wantNr)
		}
	}
}

func TestToString(t *testing.T) {
	a, err := New("0123456789")
	if err != nil {
		t.Fatalf("New(0-9) returned unexpected error %v", err)
	}
	for _, test := range []struct {
		nr         uint64
		wantString string
	}{
		{
			nr:         12,
			wantString: "12",
		},
		{
			nr:         0,
			wantString: "0",
		},
		{
			nr:         123456789,
			wantString: "123456789",
		},
		{
			nr:         9876543210,
			wantString: "9876543210",
		},
	} {
		if gotString := a.ToString(test.nr); gotString != test.wantString {
			t.Errorf("a.ToString(%v) = %v, want %v", test.nr, gotString, test.wantString)
		}
	}
}

func TestDuplicatesInTokens(t *testing.T) {
	for _, test := range []struct {
		tokens    string
		wantError bool
	}{
		{
			tokens:    "asdfghjkl",
			wantError: false,
		},
		{
			tokens:    "asdfghjkal",
			wantError: true,
		},
	} {
		_, err := New(test.tokens)
		gotError := err != nil
		if gotError != test.wantError {
			t.Errorf("New(%q): error=%v, want=%v", test.tokens, gotError, test.wantError)
		}
	}
}

func TestLargeNumbers(t *testing.T) {
	a, err := New("0123456789ABCDEFGHKLMNPQRSTUVWXYZ")
	if err != nil {
		t.Fatalf("New(0-Z) returned unexpected error %v", err)
	}
	for _, n := range []uint64{
		123456789,
		1234567890,
		123456789000,
		1234567890000,
		12345678900000,
		123456789000000,
		1234567890000000,
		12345678900000000,
		1234567890000000000,
	} {
		s := a.ToString(n)
		u, err := a.ToNr(s)
		if err != nil {
			t.Fatalf("ToNr(%q) = _,%q, need nil error", s, err.Error())
		}
		if u != n {
			t.Errorf("a.ToString(%v) = %q, but a.ToNr(%q) = %v", n, s, s, u)
		}
	}
}
