package luno

import (
	"testing"
	"time"
)

func TestMakeURLValues(t *testing.T) {

	type S string

	type Req struct {
		S   string    `url:"s"`
		I   int       `url:"i"`
		I64 int64     `url:"i64"`
		F32 float32   `url:"f32"`
		F64 float64   `url:"f64"`
		B   bool      `url:"b"`
		ABy []byte    `url:"aby"`
		TS  S         `url:"ts"`
		T   time.Time `url:"t"` // implements fmt.Stringer
	}

	r := Req{
		S:   "foo",
		I:   42,
		I64: 42,
		F32: 42.42,
		F64: 42.42,
		B:   true,
		ABy: []byte("foo"),
		TS:  S("foo"),
		T:   time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	v, err := makeURLValues(&r)
	if err != nil {
		t.Errorf("Expected success, got %v", err)
		return
	}
	exp := "aby=foo&b=true&f32=42.4200&f64=42.4200&i=42&i64=42&s=foo&t=2018-01-01+00%3A00%3A00+%2B0000+UTC&ts=foo"
	act := v.Encode()
	if act != exp {
		t.Errorf("Expected %q, got %q", exp, act)
		return
	}
}
