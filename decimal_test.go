package luno_test

import (
	"testing"

	luno "github.com/luno/luno-go"
)

func TestDecimalString(t *testing.T) {
	type testCase struct {
		d   luno.Decimal
		exp string
	}

	testCases := []testCase{
		testCase{
			d:   luno.NewDecimal(0, 0),
			exp: "0",
		},
		testCase{
			d:   luno.NewDecimal(1, 0),
			exp: "1",
		},
		testCase{
			d:   luno.NewDecimal(1.2, 1),
			exp: "1.2",
		},
		testCase{
			d:   luno.NewDecimal(-1.2, 1),
			exp: "-1.2",
		},
		testCase{
			d:   luno.NewDecimal(1.2, 0),
			exp: "1",
		},
		testCase{
			d:   luno.NewDecimal(1.12345678, 8),
			exp: "1.12345678",
		},
		testCase{
			d:   luno.NewDecimal(-1.12345678, 8),
			exp: "-1.12345678",
		},
		testCase{
			d:   luno.NewDecimal(1.123456789, 8),
			exp: "1.12345679", // round towards +inf
		},
		testCase{
			d:   luno.NewDecimal(-1.123456789, 8),
			exp: "-1.12345679", // round towards -inf
		},
		testCase{
			d:   luno.NewDecimal(1.123456782, 8),
			exp: "1.12345678", // round towards 0
		},
		testCase{
			d:   luno.NewDecimal(-1.123456782, 8),
			exp: "-1.12345678", // round towards 0
		},
	}

	for _, test := range testCases {
		act := test.d.String()
		if act != test.exp {
			t.Errorf("Expected string representation of %v to be %q, not %q",
				test.d, test.exp, act)
		}
	}
}
