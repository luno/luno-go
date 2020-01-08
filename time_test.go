package luno_test

import (
	"strconv"
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
		testCase{
			err: false,
		},
		testCase{
			in:  []byte{},
			err: false,
		},
		testCase{
			in:  []byte("abc"),
			err: true,
		},
		testCase{
			in:  []byte("123456"),
			exp: luno.Time(time.Unix(0, 123456e6)),
		},
		testCase{
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
		mct, err := act.MarshalJSON()
		if err != nil {
			t.Errorf("Expected marshalling unmarshalled %q to succeed", string(test.in))
		}
		if string(mct) != string(test.in) {
			t.Errorf("Expected marshalling unmarshalled %q to have original value %q", string(mct), string(test.in))
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
		testCase{
			in: luno.Time{},
		},
		testCase{
			in:  luno.Time(now),
			exp: strconv.FormatInt(now.UnixNano()/1e6, 10),
		},
		testCase{
			in:  luno.Time(time.Date(2006, 1, 2, 3, 4, 5, 999, time.UTC)),
			exp: "1136171045000",
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
