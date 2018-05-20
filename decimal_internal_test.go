package luno

import "testing"

func TestDecimalUnmarshalJSON(t *testing.T) {
	type testCase struct {
		b         string
		err       bool
		f         float64
		precision int
	}

	testCases := []testCase{
		testCase{
			b:   `""`,
			err: true,
		},
		testCase{
			b:   `"abc"`,
			err: true,
		},
		testCase{
			b:         `"0."`,
			err:       false,
			f:         0,
			precision: 0,
		},
		testCase{
			b:         `"0"`,
			err:       false,
			f:         0,
			precision: 0,
		},
		testCase{
			b:         `"0.1"`,
			err:       false,
			f:         0.1,
			precision: 1,
		},
		testCase{
			b:         `"-100.1"`,
			err:       false,
			f:         -100.1,
			precision: 1,
		},
		testCase{
			b:         `"1.12345678"`,
			err:       false,
			f:         1.12345678,
			precision: 8,
		},
		testCase{
			b:         `"1.123456789"`,
			err:       false,
			f:         1.123456789,
			precision: 9,
		},
		testCase{
			b:   `"1e8"`,
			err: true,
		},
		testCase{
			b:   `"1.1e8"`,
			err: true,
		},
		testCase{
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
		if d.f != test.f {
			t.Errorf("Expected %q to unmarshal as %f, got %f",
				string(test.b), test.f, d.f)
		}
		if d.precision != test.precision {
			t.Errorf("Expected %q to unmarshal with precision %d, got %d",
				string(test.b), test.precision, d.precision)
		}
	}
}
