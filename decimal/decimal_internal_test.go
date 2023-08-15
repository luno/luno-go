package decimal

import (
	"math/big"
	"testing"
)

func TestDecimalUnmarshalJSON(t *testing.T) {
	type testCase struct {
		b     string
		err   bool
		i     *big.Int
		scale int
	}

	testCases := []testCase{
		{
			b:   `""`,
			err: true,
		},
		{
			b:   `"abc"`,
			err: true,
		},
		{
			b:     `"0."`,
			err:   false,
			i:     big.NewInt(0),
			scale: 0,
		},
		{
			b:     `"0"`,
			err:   false,
			i:     big.NewInt(0),
			scale: 0,
		},
		{
			b:     `"0.1"`,
			err:   false,
			i:     big.NewInt(1),
			scale: 1,
		},
		{
			b:     `"-100.1"`,
			err:   false,
			i:     big.NewInt(-1001),
			scale: 1,
		},
		{
			b:     `"1.12345678"`,
			err:   false,
			i:     big.NewInt(112345678),
			scale: 8,
		},
		{
			b:     `"1.123456789"`,
			err:   false,
			i:     big.NewInt(1123456789),
			scale: 9,
		},
		{
			b:   `"1e8"`,
			err: true,
		},
		{
			b:   `"1.1e8"`,
			err: true,
		},
		{
			b:   `"1.1234e2"`,
			err: true,
		},
	}

	var d Decimal
	for _, test := range testCases {
		err := (&d).UnmarshalJSON([]byte(test.b))
		if err != nil {
			if !test.err {
				t.Errorf("Expected unmarhsalling %q to succeed, got %v",
					string(test.b), err)
			}
			continue
		} else if test.err {
			t.Errorf("Expected unmarhsalling %q to fail", string(test.b))
			continue
		}
		if d.i.Cmp(test.i) != 0 {
			t.Errorf("Expected %q to unmarshal as %s, got %s",
				string(test.b), test.i, d.i)
		}
		if d.scale != test.scale {
			t.Errorf("Expected %q to unmarshal with scale %d, got %d",
				string(test.b), test.scale, d.scale)
		}
	}
}

func TestDecimalToScale(t *testing.T) {
	type testCase struct {
		d     Decimal
		scale int
		exp   *big.Int
	}

	testCases := []testCase{
		{
			d:     New(big.NewInt(0), 0),
			scale: 0,
			exp:   big.NewInt(0),
		},
		{
			d:     New(big.NewInt(1), 0),
			scale: 1,
			exp:   big.NewInt(10),
		},
		{
			d:     New(big.NewInt(-1), 0),
			scale: 1,
			exp:   big.NewInt(-10),
		},
		{
			d:     New(big.NewInt(12344), 4),
			scale: 3,
			exp:   big.NewInt(1234),
		},
		{
			d:     New(big.NewInt(-12344), 4),
			scale: 3,
			exp:   big.NewInt(-1234),
		},
		{
			d:     New(big.NewInt(12349), 4),
			scale: 3,
			exp:   big.NewInt(1234),
		},
		{
			d:     New(big.NewInt(-12349), 4),
			scale: 3,
			exp:   big.NewInt(-1234),
		},
	}

	for _, test := range testCases {
		d2 := test.d.ToScale(test.scale)
		act := d2.i
		if act.Cmp(test.exp) != 0 {
			t.Errorf("Expected %s scaled to %d to be %s, not %s",
				test.d, test.scale, test.exp, act)
		}
	}
}
