package luno

import (
	"errors"
	"strconv"
	"strings"
)

type Decimal struct {
	f         float64
	precision int
}

func NewDecimal(f float64, precision int) Decimal {
	return Decimal{f: f, precision: precision}
}

func (d Decimal) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

func (d *Decimal) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if strings.IndexByte(s, 'e') > -1 {
		return errors.New("luno: unsupported decimal notation")
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*d = NewDecimal(f, getPrecision(s))
	return nil
}

func (d Decimal) String() string {
	return strconv.FormatFloat(d.f, 'f', d.precision, 64)
}

func (d Decimal) Float64() float64 {
	return d.f
}

func getPrecision(s string) int {
	i := strings.IndexByte(s, '.')
	if i == -1 {
		return 0
	}
	return len(s) - i - 1
}
