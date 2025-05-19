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
			in:  []byte("null"),
			exp: luno.Time{},
		},
		{
			in:  []byte{},
			exp: luno.Time{},
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
		{
			in:  []byte(`"2006-01-02T15:04:05Z"`),
			exp: luno.Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
		},
		{
			in:  []byte(`"2006-01-02 15:04:05 +0000 UTC"`),
			exp: luno.Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
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

	date := time.Date(2006, 1, 2, 3, 4, 5, 0, time.UTC)

	testCases := []testCase{
		{
			in:  luno.Time{},
			exp: "null",
		},
		{
			in:  luno.Time(date),
			exp: `"2006-01-02T03:04:05Z"`,
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
