package luno

import (
	"testing"
	"time"
)

func TestMakeURLValues(t *testing.T) {

	type S string

	type Req struct {
		S   string  `url:"s"`
		I   int     `url:"i"`
		I64 int64   `url:"i64"`
		F32 float32 `url:"f32"`
		F64 float64 `url:"f64"`
		B   bool    `url:"b"`
		ABy []byte  `url:"aby"`
		TS  S       `url:"ts"`
		T   Time    `url:"t"` // implements QueryValuer
	}

	tests := []struct {
		name     string
		r        Req
		expected string
	}{
		{
			name: "valid time",
			r: Req{
				S:   "foo",
				I:   42,
				I64: 42,
				F32: 42.42,
				F64: 42.42,
				B:   true,
				ABy: []byte("foo"),
				TS:  S("foo"),
				T:   Time(time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			expected: "aby=foo&b=true&f32=42.4200&f64=42.4200&i=42&i64=42&s=foo&t=1514764800000&ts=foo",
		},
		{
			name: "zero time",
			r: Req{
				S:   "foo",
				I:   42,
				I64: 42,
				F32: 42.42,
				F64: 42.42,
				B:   true,
				ABy: []byte("foo"),
				TS:  S("foo"),
				T:   Time(time.Time{}),
			},
			expected: "aby=foo&b=true&f32=42.4200&f64=42.4200&i=42&i64=42&s=foo&t=&ts=foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := makeURLValues(&(tt.r))
			if err != nil {
				t.Errorf("Expected success, got %v", err)
				return
			}
			act := v.Encode()
			if act != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, act)
				return
			}
		})
	}
}
