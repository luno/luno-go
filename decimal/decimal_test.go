package decimal_test

import (
	"math/big"
	"testing"

	"github.com/luno/luno-go/decimal"
)

func TestFloat64(t *testing.T) {
	type testCase struct {
		a   string
		b   string
		exp float64
	}

	testCases := []testCase{
		{
			a:   "0.17800000",
			b:   "6750.00000000",
			exp: 1201.5,
		},
		{
			a:   "0.178000001",
			b:   "6750.000000001",
			exp: 1201.500006750178,
		},
	}

	for _, test := range testCases {
		a, _ := decimal.NewFromString(test.a)
		b, _ := decimal.NewFromString(test.b)

		act := a.Mul(b).Float64()

		if act != test.exp {
			t.Errorf("Expected %s * %s to stringify as %g, got %g", test.a, test.b, test.exp, act)
		}
	}
}

func TestNewFromInt64(t *testing.T) {
	type testCase struct {
		i   int64
		exp string
	}

	testCases := []testCase{
		{
			i:   0,
			exp: "0",
		},
		{
			i:   1,
			exp: "1",
		},
		{
			i:   -1,
			exp: "-1",
		},
		{
			i:   1231231,
			exp: "1231231",
		},
	}

	for _, test := range testCases {
		act := decimal.NewFromInt64(test.i).String()
		if act != test.exp {
			t.Errorf("Expected %d to stringify as %q, got %q",
				test.i, test.exp, act)
		}
	}
}

func TestNewFromFloat64(t *testing.T) {
	type testCase struct {
		f     float64
		scale int
		exp   string
	}

	testCases := []testCase{
		{
			f:     0,
			scale: 0,
			exp:   "0",
		},
		{
			f:     1,
			scale: 0,
			exp:   "1",
		},
		{
			f:     1.12345678,
			scale: 0,
			exp:   "1",
		},
		{
			f:     1.12345678,
			scale: 8,
			exp:   "1.12345678",
		},
		{
			f:     -1.12345678,
			scale: 4,
			exp:   "-1.1234",
		},
	}

	for _, test := range testCases {
		act := decimal.NewFromFloat64(test.f, test.scale).String()
		if act != test.exp {
			t.Errorf("Expected %f (scale %d) to stringify as %q, got %q",
				test.f, test.scale, test.exp, act)
		}
	}
}

func TestNewFromString(t *testing.T) {
	type testCase struct {
		s   string
		err bool
	}

	testCases := []testCase{
		{
			s:   "",
			err: true,
		},
		{
			s:   "abc",
			err: true,
		},
		{
			s:   "1e8",
			err: true,
		},
		{s: "0"},
		{s: "1"},
		{s: "-1.2"},
		{s: "1.12345678"},
		{s: "1.123456789"},
	}

	for _, test := range testCases {
		d, err := decimal.NewFromString(test.s)
		if err != nil {
			if !test.err {
				t.Errorf("Expected %q to succeed, got %v", test.s, err)
			}
			continue
		} else if test.err {
			t.Errorf("Expected %q to fail", test.s)
			continue
		}
		act := d.String()
		if act != test.s {
			t.Errorf("%q failed to stringify back to itself, got %q",
				test.s, act)
		}
	}
}

func TestZero(t *testing.T) {
	d := decimal.Zero()
	act := d.String()
	exp := "0"
	if act != exp {
		t.Errorf("Expected Zero() to return %q, got %q", exp, act)
	}
}

func TestDecimalString(t *testing.T) {
	type testCase struct {
		d   decimal.Decimal
		exp string
	}

	testCases := []testCase{
		{
			d:   decimal.New(big.NewInt(0), 0),
			exp: "0",
		},
		{
			d:   decimal.New(big.NewInt(1), 0),
			exp: "1",
		},
		{
			d:   decimal.New(big.NewInt(12), 1),
			exp: "1.2",
		},
		{
			d:   decimal.New(big.NewInt(-12), 1),
			exp: "-1.2",
		},
		{
			d:   decimal.New(big.NewInt(12), 0),
			exp: "12",
		},
		{
			d:   decimal.New(big.NewInt(112345678), 8),
			exp: "1.12345678",
		},
		{
			d:   decimal.New(big.NewInt(-112345678), 8),
			exp: "-1.12345678",
		},
		{
			d:   decimal.New(big.NewInt(1123456789), 9),
			exp: "1.123456789",
		},
		{
			d:   decimal.New(big.NewInt(-1123456789), 9),
			exp: "-1.123456789",
		},
		{
			d:   decimal.New(big.NewInt(1123456782), 9),
			exp: "1.123456782",
		},
		{
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
		{
			d:   decimal.New(big.NewInt(0), 0),
			exp: 0,
		},
		{
			d:   decimal.New(big.NewInt(1), 0),
			exp: 1,
		},
		{
			d:   decimal.New(big.NewInt(-1), 0),
			exp: -1,
		},
		{
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
		{
			d1:  decimal.New(big.NewInt(0), 0),
			d2:  decimal.New(big.NewInt(0), 0),
			exp: 0,
		},
		{
			d1:  decimal.New(big.NewInt(1), 0),
			d2:  decimal.New(big.NewInt(0), 0),
			exp: 1,
		},
		{
			d1:  decimal.New(big.NewInt(-1), 0),
			d2:  decimal.New(big.NewInt(0), 0),
			exp: -1,
		},
		{
			d1:  decimal.New(big.NewInt(100), 1),
			d2:  decimal.New(big.NewInt(100), 3),
			exp: 1,
		},
		{
			d1:  decimal.New(big.NewInt(100), 3),
			d2:  decimal.New(big.NewInt(100), 1),
			exp: -1,
		},
		{
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
		{
			d:   decimal.New(big.NewInt(0), 0),
			exp: decimal.New(big.NewInt(0), 0),
		},
		{
			d:   decimal.New(big.NewInt(1), 0),
			exp: decimal.New(big.NewInt(-1), 0),
		},
		{
			d:   decimal.New(big.NewInt(-1), 0),
			exp: decimal.New(big.NewInt(1), 0),
		},
		{
			d:   decimal.New(big.NewInt(5), 1),
			exp: decimal.New(big.NewInt(-5), 1),
		},
		{
			d:   decimal.New(big.NewInt(-5), 1),
			exp: decimal.New(big.NewInt(5), 1),
		},
		{
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
		{
			d1:        decimal.New(big.NewInt(0), 0),
			d2:        decimal.New(big.NewInt(0), 0),
			exp:       decimal.New(big.NewInt(0), 0),
			expString: "0",
		},
		{
			d1:        decimal.New(big.NewInt(1), 0),
			d2:        decimal.New(big.NewInt(0), 0),
			exp:       decimal.New(big.NewInt(1), 0),
			expString: "1",
		},
		{
			d1:        decimal.New(big.NewInt(-1), 0),
			d2:        decimal.New(big.NewInt(0), 0),
			exp:       decimal.New(big.NewInt(-1), 0),
			expString: "-1",
		},
		{
			d1:        decimal.New(big.NewInt(1), 0),
			d2:        decimal.New(big.NewInt(112345678), 8),
			exp:       decimal.New(big.NewInt(212345678), 8),
			expString: "2.12345678",
		},
		{
			d1:        decimal.New(big.NewInt(-1123), 3),
			d2:        decimal.New(big.NewInt(-1123), 3),
			exp:       decimal.New(big.NewInt(-2246), 3),
			expString: "-2.246",
		},
		{
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
		{
			d1:        decimal.New(big.NewInt(0), 0),
			d2:        decimal.New(big.NewInt(0), 0),
			exp:       decimal.New(big.NewInt(0), 0),
			expString: "0",
		},
		{
			d1:        decimal.New(big.NewInt(1), 0),
			d2:        decimal.New(big.NewInt(0), 0),
			exp:       decimal.New(big.NewInt(1), 0),
			expString: "1",
		},
		{
			d1:        decimal.New(big.NewInt(-1), 0),
			d2:        decimal.New(big.NewInt(0), 0),
			exp:       decimal.New(big.NewInt(-1), 0),
			expString: "-1",
		},
		{
			d1:        decimal.New(big.NewInt(1), 0),
			d2:        decimal.New(big.NewInt(112345678), 8),
			exp:       decimal.New(big.NewInt(-12345678), 8),
			expString: "-0.12345678",
		},
		{
			d1:        decimal.New(big.NewInt(-1123), 3),
			d2:        decimal.New(big.NewInt(-1123), 3),
			exp:       decimal.New(big.NewInt(0), 3),
			expString: "0.000",
		},
		{
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

func TestDecimalMulInt64(t *testing.T) {
	type testCase struct {
		d   decimal.Decimal
		y   int64
		exp string
	}

	testCases := []testCase{
		{
			d:   decimal.New(big.NewInt(0), 0),
			y:   0,
			exp: "0",
		},
		{
			d:   decimal.New(big.NewInt(100), 0),
			y:   100,
			exp: "10000",
		},
		{
			d:   decimal.New(big.NewInt(-100), 0),
			y:   100,
			exp: "-10000",
		},
		{
			d:   decimal.New(big.NewInt(100), 8),
			y:   100,
			exp: "0.00010000",
		},
		{
			d:   decimal.New(big.NewInt(-100), 8),
			y:   100,
			exp: "-0.00010000",
		},
		{
			d:   decimal.New(big.NewInt(17823), 4),
			y:   124,
			exp: "221.0052",
		},
	}

	for _, test := range testCases {
		act := test.d.MulInt64(test.y).String()
		if act != test.exp {
			t.Errorf("Expected %s * %d to be %q, got %q",
				test.d, test.y, test.exp, act)
		}
	}
}

func TestDecimalDivInt64(t *testing.T) {
	type testCase struct {
		d   decimal.Decimal
		y   int64
		exp string
	}

	testCases := []testCase{
		{
			d:   decimal.New(big.NewInt(1), 0),
			y:   1,
			exp: "1",
		},
		{
			d:   decimal.New(big.NewInt(-100), 0),
			y:   100,
			exp: "-1",
		},
		{
			d:   decimal.New(big.NewInt(100), 8),
			y:   100,
			exp: "0.00000001",
		},
		{
			d:   decimal.New(big.NewInt(-100), 8),
			y:   100,
			exp: "-0.00000001",
		},
		{
			d:   decimal.New(big.NewInt(1), 4),
			y:   10,
			exp: "0.0000",
		},
		{
			d:   decimal.New(big.NewInt(17823), 4),
			y:   124,
			exp: "0.0143",
		},
	}

	for _, test := range testCases {
		act := test.d.DivInt64(test.y).String()
		if act != test.exp {
			t.Errorf("Expected %s / %d to be %q, got %q",
				test.d, test.y, test.exp, act)
		}
	}
}

func TestDecimalMul(t *testing.T) {
	type testCase struct {
		d   decimal.Decimal
		y   decimal.Decimal
		exp string
	}

	testCases := []testCase{
		{
			d:   decimal.Decimal{},
			y:   decimal.Decimal{},
			exp: "0",
		},
		{
			d:   decimal.New(big.NewInt(1), 0),
			y:   decimal.New(big.NewInt(1), 0),
			exp: "1",
		},
		{
			d:   decimal.New(big.NewInt(1), 0),
			y:   decimal.New(big.NewInt(-1), 0),
			exp: "-1",
		},
		{
			d:   decimal.New(big.NewInt(1), 3),
			y:   decimal.New(big.NewInt(1), 3),
			exp: "0.000001",
		},
		{
			d:   decimal.New(big.NewInt(1234), 3),
			y:   decimal.New(big.NewInt(5678), 8),
			exp: "0.00007006652",
		},
		{
			d:   decimal.New(big.NewInt(100), 2),
			y:   decimal.New(big.NewInt(100), 2),
			exp: "1.0000",
		},
	}

	for _, test := range testCases {
		act := test.d.Mul(test.y).String()
		if act != test.exp {
			t.Errorf("Expected %s * %s to be %q, got %q",
				test.d, test.y, test.exp, act)
		}
	}
}

func TestDecimalDivNoPanic(t *testing.T) {
	type testCase struct {
		d     decimal.Decimal
		y     decimal.Decimal
		scale int
		exp   string
	}

	testCases := []testCase{
		{
			d:     decimal.New(big.NewInt(1), 0),
			y:     decimal.New(big.NewInt(1), 0),
			scale: 0,
			exp:   "1",
		},
		{
			d:     decimal.New(big.NewInt(1), 0),
			y:     decimal.New(big.NewInt(10), 0),
			scale: 2,
			exp:   "0.10",
		},
		{
			d:     decimal.New(big.NewInt(1), 0),
			y:     decimal.New(big.NewInt(-10), 0),
			scale: 2,
			exp:   "-0.10",
		},
		{
			d:     decimal.New(big.NewInt(5000), 4),
			y:     decimal.New(big.NewInt(15000), 4),
			scale: 2,
			exp:   "0.33",
		},
		{
			d:     decimal.New(big.NewInt(5000), 4),
			y:     decimal.New(big.NewInt(15000), 4),
			scale: 10,
			exp:   "0.3333333333",
		},
		{
			d:     decimal.New(big.NewInt(1234), 2),
			y:     decimal.New(big.NewInt(5678), 8),
			scale: 5,
			exp:   "217330.04579",
		},
	}

	for _, test := range testCases {
		act := test.d.Div(test.y, test.scale).String()
		if act != test.exp {
			t.Errorf("Expected %s / %s to be %q, got %q",
				test.d, test.y, test.exp, act)
		}
	}
}

func testDecimalDivPanic(t *testing.T, d, y decimal.Decimal) {
	defer func() {
		if e := recover(); e == nil {
			t.Errorf("Expected %s / %s to panic", d, y)
		}
	}()

	d.Div(y, 0)
}

func TestDecimalDivPanic(t *testing.T) {
	type testCase struct {
		d decimal.Decimal
		y decimal.Decimal
	}

	testCases := []testCase{
		{
			d: decimal.Decimal{},
			y: decimal.Decimal{},
		},
		{
			d: decimal.New(big.NewInt(1), 0),
			y: decimal.Decimal{},
		},
		{
			d: decimal.New(big.NewInt(1), 0),
			y: decimal.New(big.NewInt(0), 0),
		},
	}

	for _, test := range testCases {
		testDecimalDivPanic(t, test.d, test.y)
	}
}
