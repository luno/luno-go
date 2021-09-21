package luno

import (
	"testing"
	"time"

	"github.com/luno/luno-go/decimal"
)

func TestMakeURLValues(t *testing.T) {

	type S string

	type Req struct {
		S   string          `url:"s"`
		I   int             `url:"i"`
		I64 int64           `url:"i64"`
		F32 float32         `url:"f32"`
		F64 float64         `url:"f64"`
		B   bool            `url:"b"`
		ASt []string        `url:"ast"`
		TS  S               `url:"ts"`
		Amt decimal.Decimal `url:"amt"`
		T   Time            `url:"t"` // implements QueryValuer
	}

	tests := []struct {
		name     string
		r        *Req
		expected string
	}{
		{
			name: "nil",
			r: nil,
		},
		{
			name: "valid time",
			r: &Req{
				S:   "foo",
				I:   42,
				I64: 42,
				F32: 42.42,
				F64: 42.42,
				B:   true,
				ASt: []string{"foo", "bar"},
				TS:  S("foo"),
				Amt: decimal.NewFromFloat64(2.1, 1),
				T:   Time(time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			expected: "amt=2.1&ast=foo&ast=bar&b=true&f32=42.4200&f64=42.4200&i=42&i64=42&s=foo&t=1514764800000&ts=foo",
		},
		{
			name: "zero time",
			r: &Req{
				S:   "foo",
				I:   42,
				I64: 42,
				F32: 42.42,
				F64: 42.42,
				B:   true,
				ASt: []string{"foo", "bar"},
				TS:  S("foo"),
				Amt: decimal.NewFromFloat64(0.1, 1),
				T:   Time(time.Time{}),
			},
			expected: "amt=0.1&ast=foo&ast=bar&b=true&f32=42.4200&f64=42.4200&i=42&i64=42&s=foo&t=&ts=foo",
		},
		{
			name:"valid amount",
			r: &Req{
				S:   "foo",
				I:   42,
				I64: 42,
				F32: 42.42,
				F64: 42.42,
				B:   true,
				ASt: []string{"foo", "bar"},
				TS:  S("foo"),
				Amt: decimal.NewFromFloat64(0.0001, 10),
				T:   Time(time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			expected:"amt=0.0001000000&ast=foo&ast=bar&b=true&f32=42.4200&f64=42.4200&i=42&i64=42&s=foo&t=1514764800000&ts=foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := makeURLValues(tt.r).Encode()
			if act != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, act)
				return
			}
		})
	}
}

func TestNilURLValues(t *testing.T) {
	act := makeURLValues(nil).Encode()
	if act != "" {
		t.Errorf("Expected empty url, got %s", act)
	}
}
