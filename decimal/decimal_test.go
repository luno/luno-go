package decimal_test

import (
	"math/big"
	"testing"

	"github.com/luno/luno-go/decimal"
)

func TestDecimalString(t *testing.T) {
	type testCase struct {
		d   decimal.Decimal
		exp string
	}

	testCases := []testCase{
		testCase{
			d:   decimal.New(big.NewInt(0), 0),
			exp: "0",
		},
		testCase{
			d:   decimal.New(big.NewInt(1), 0),
			exp: "1",
		},
		testCase{
			d:   decimal.New(big.NewInt(12), 1),
			exp: "1.2",
		},
		testCase{
			d:   decimal.New(big.NewInt(-12), 1),
			exp: "-1.2",
		},
		testCase{
			d:   decimal.New(big.NewInt(12), 0),
			exp: "12",
		},
		testCase{
			d:   decimal.New(big.NewInt(112345678), 8),
			exp: "1.12345678",
		},
		testCase{
			d:   decimal.New(big.NewInt(-112345678), 8),
			exp: "-1.12345678",
		},
		testCase{
			d:   decimal.New(big.NewInt(1123456789), 9),
			exp: "1.123456789",
		},
		testCase{
			d:   decimal.New(big.NewInt(-1123456789), 9),
			exp: "-1.123456789",
		},
		testCase{
			d:   decimal.New(big.NewInt(1123456782), 9),
			exp: "1.123456782",
		},
		testCase{
			d:   decimal.New(big.NewInt(-1123456782), 9),
			exp: "-1.123456782",
		},
	}

	for _, test := range testCases {
		act := test.d.String()
		if act != test.exp {
			t.Errorf("Expected string of %v to be %q, not %q",
				test.d, test.exp, act)
		}
	}
}

func TestDecimalSign(t *testing.T) {
	type testCase struct {
		d   decimal.Decimal
		exp int
	}

	testCases := []testCase{
		testCase{
			d:   decimal.New(big.NewInt(0), 0),
			exp: 0,
		},
		testCase{
			d:   decimal.New(big.NewInt(1), 0),
			exp: 1,
		},
		testCase{
			d:   decimal.New(big.NewInt(-1), 0),
			exp: -1,
		},
		testCase{
			d:   decimal.New(big.NewInt(12345678), 8),
			exp: 1,
		},
	}

	for _, test := range testCases {
		act := test.d.Sign()
		if act != test.exp {
			t.Errorf("Expected the sign of %s to be %d, got %d",
				test.d, test.exp, act)
		}
	}
}

func TestDecimalCmp(t *testing.T) {
	type testCase struct {
		d1  decimal.Decimal
		d2  decimal.Decimal
		exp int
	}

	testCases := []testCase{
		testCase{
			d1:  decimal.New(big.NewInt(0), 0),
			d2:  decimal.New(big.NewInt(0), 0),
			exp: 0,
		},
		testCase{
			d1:  decimal.New(big.NewInt(1), 0),
			d2:  decimal.New(big.NewInt(0), 0),
			exp: 1,
		},
		testCase{
			d1:  decimal.New(big.NewInt(-1), 0),
			d2:  decimal.New(big.NewInt(0), 0),
			exp: -1,
		},
		testCase{
			d1:  decimal.New(big.NewInt(100), 1),
			d2:  decimal.New(big.NewInt(100), 3),
			exp: 1,
		},
		testCase{
			d1:  decimal.New(big.NewInt(100), 3),
			d2:  decimal.New(big.NewInt(100), 1),
			exp: -1,
		},
		testCase{
			d1:  decimal.New(big.NewInt(100), 2),
			d2:  decimal.New(big.NewInt(1), 0),
			exp: 0,
		},
	}

	for _, test := range testCases {
		act := test.d1.Cmp(test.d2)
		if act != test.exp {
			t.Errorf("Expected %d when comparing %s to %s, got %d",
				test.exp, test.d2, test.d1, act)
		}
	}
}

func TestDecimalNeg(t *testing.T) {
	type testCase struct {
		d   decimal.Decimal
		exp decimal.Decimal
	}

	testCases := []testCase{
		testCase{
			d:   decimal.New(big.NewInt(0), 0),
			exp: decimal.New(big.NewInt(0), 0),
		},
		testCase{
			d:   decimal.New(big.NewInt(1), 0),
			exp: decimal.New(big.NewInt(-1), 0),
		},
		testCase{
			d:   decimal.New(big.NewInt(-1), 0),
			exp: decimal.New(big.NewInt(1), 0),
		},
		testCase{
			d:   decimal.New(big.NewInt(5), 1),
			exp: decimal.New(big.NewInt(-5), 1),
		},
		testCase{
			d:   decimal.New(big.NewInt(-5), 1),
			exp: decimal.New(big.NewInt(5), 1),
		},
		testCase{
			d:   decimal.New(big.NewInt(12345678), 8),
			exp: decimal.New(big.NewInt(-12345678), 8),
		},
	}

	for _, test := range testCases {
		act := test.d.Neg()
		if act.Cmp(test.exp) != 0 {
			t.Errorf("Expected the negation of %s to be %s, got %s",
				test.d, test.exp, act)
		}
	}
}

func TestDecimalAdd(t *testing.T) {
	type testCase struct {
		d1        decimal.Decimal
		d2        decimal.Decimal
		exp       decimal.Decimal
		expString string
	}

	testCases := []testCase{
		testCase{
			d1:        decimal.New(big.NewInt(0), 0),
			d2:        decimal.New(big.NewInt(0), 0),
			exp:       decimal.New(big.NewInt(0), 0),
			expString: "0",
		},
		testCase{
			d1:        decimal.New(big.NewInt(1), 0),
			d2:        decimal.New(big.NewInt(0), 0),
			exp:       decimal.New(big.NewInt(1), 0),
			expString: "1",
		},
		testCase{
			d1:        decimal.New(big.NewInt(-1), 0),
			d2:        decimal.New(big.NewInt(0), 0),
			exp:       decimal.New(big.NewInt(-1), 0),
			expString: "-1",
		},
		testCase{
			d1:        decimal.New(big.NewInt(1), 0),
			d2:        decimal.New(big.NewInt(112345678), 8),
			exp:       decimal.New(big.NewInt(212345678), 8),
			expString: "2.12345678",
		},
		testCase{
			d1:        decimal.New(big.NewInt(-1123), 3),
			d2:        decimal.New(big.NewInt(-1123), 3),
			exp:       decimal.New(big.NewInt(-2246), 3),
			expString: "-2.246",
		},
		testCase{
			d1:        decimal.New(big.NewInt(112345678), 8),
			d2:        decimal.New(big.NewInt(112345678), 8),
			exp:       decimal.New(big.NewInt(224691356), 8),
			expString: "2.24691356",
		},
	}

	for _, test := range testCases {
		act := test.d1.Add(test.d2)
		if act.Cmp(test.exp) != 0 {
			t.Errorf("Expected the sum of %s and %s to be %s, got %s",
				test.d1, test.d2, test.exp, act)
			continue
		}

		actString := act.String()
		if actString != test.expString {
			t.Errorf("Expected the sum of %s and %s to stringify as %q, got %q",
				test.d1, test.d2, test.expString, actString)
		}
	}
}

func TestDecimalSub(t *testing.T) {
	type testCase struct {
		d1        decimal.Decimal
		d2        decimal.Decimal
		exp       decimal.Decimal
		expString string
	}

	testCases := []testCase{
		testCase{
			d1:        decimal.New(big.NewInt(0), 0),
			d2:        decimal.New(big.NewInt(0), 0),
			exp:       decimal.New(big.NewInt(0), 0),
			expString: "0",
		},
		testCase{
			d1:        decimal.New(big.NewInt(1), 0),
			d2:        decimal.New(big.NewInt(0), 0),
			exp:       decimal.New(big.NewInt(1), 0),
			expString: "1",
		},
		testCase{
			d1:        decimal.New(big.NewInt(-1), 0),
			d2:        decimal.New(big.NewInt(0), 0),
			exp:       decimal.New(big.NewInt(-1), 0),
			expString: "-1",
		},
		testCase{
			d1:        decimal.New(big.NewInt(1), 0),
			d2:        decimal.New(big.NewInt(112345678), 8),
			exp:       decimal.New(big.NewInt(-12345678), 8),
			expString: "-0.12345678",
		},
		testCase{
			d1:        decimal.New(big.NewInt(-1123), 3),
			d2:        decimal.New(big.NewInt(-1123), 3),
			exp:       decimal.New(big.NewInt(0), 3),
			expString: "0.000",
		},
		testCase{
			d1:        decimal.New(big.NewInt(112345678), 8),
			d2:        decimal.New(big.NewInt(112345678), 8),
			exp:       decimal.New(big.NewInt(0), 8),
			expString: "0.00000000",
		},
	}

	for _, test := range testCases {
		act := test.d1.Sub(test.d2)
		if act.Cmp(test.exp) != 0 {
			t.Errorf("Expected %s - %s to be %s, got %s",
				test.d1, test.d2, test.exp, act)
			continue
		}

		actString := act.String()
		if actString != test.expString {
			t.Errorf("Expected %s - %s to stringify as %q, got %q",
				test.d1, test.d2, test.expString, actString)
		}
	}
}
