package luno_test

import (
	"testing"
	"time"

	luno "github.com/luno/luno-go"
)

func TestTimeUnmarshalJSON(t *testing.T) {
	type testCase struct {
		in  []byte
		err bool
		exp luno.Time
	}

	testCases := []testCase{
		{
			err: true,
		},
		{
			in:  []byte{},
			err: true,
		},
		{
			in:  []byte("abc"),
			err: true,
		},
		{
			in:  []byte("123456"),
			exp: luno.Time(time.Unix(0, 123456e6)),
		},
		{
			in:  []byte("-123456"),
			exp: luno.Time(time.Unix(0, -123456e6)),
		},
	}

	var act luno.Time
	for _, test := range testCases {
		err := (&act).UnmarshalJSON(test.in)
		if err != nil {
			if !test.err {
				t.Errorf("Expected unmarshalling %q to fail", string(test.in))
			}
			continue
		}
		if act != test.exp {
			t.Errorf("Expected %q to unmarshal as %v, got %v",
				string(test.in), test.exp, act)
		}
	}
}

func TestTimeMarshalJSON(t *testing.T) {
	type testCase struct {
		in  luno.Time
		exp string
	}

	now := time.Now()

	testCases := []testCase{
		{
			in:  luno.Time{},
			exp: time.Time{}.String(),
		},
		{
			in:  luno.Time(now),
			exp: now.String(),
		},
		{
			in:  luno.Time(time.Date(2006, 1, 2, 3, 4, 5, 999, time.UTC)),
			exp: time.Date(2006, 1, 2, 3, 4, 5, 999, time.UTC).String(),
		},
	}

	for _, test := range testCases {
		b, err := test.in.MarshalJSON()
		if err != nil {
			t.Errorf("Expected marshalling %v to succeed", test.in)
			continue
		}
		act := string(b)
		if act != test.exp {
			t.Errorf("Expected %v to marshal as %q, got %q",
				test.in, test.exp, act)
		}
	}
}
