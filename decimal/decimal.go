package decimal

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
)

// Maintain references to common big.Ints to reduce GC pressure
var bigIntLookup = [64]*big.Int{}

func init() {
	for i := 0; i < len(bigIntLookup); i++ {
		bigIntLookup[i] = big.NewInt(int64(i))
	}
}

func cachedBigInt(i int64) *big.Int {
	if i >= 0 && i < int64(len(bigIntLookup)) {
		return bigIntLookup[i]
	}
	return big.NewInt(i)
}

// Decimal represents a decimal amount, internally stored as a big.Int.
type Decimal struct {
	i     *big.Int
	scale int
}

// New returns a new Decimal. scale specifies the units of i as an inverse power
// of 10. For example, a scale of 2 indicates that i is in units of 10^-2, and
// a value of i=25 would represent 25*10^-2=0.25. If scale=0, then the decimal
// is an integer.
func New(i *big.Int, scale int) Decimal {
	return Decimal{
		i:     i,
		scale: scale,
	}
}

func NewFromInt64(i int64) Decimal {
	return New(cachedBigInt(i), 0)
}

// NewFromFloat64 returns a new Decimal with the given scale. The value is
// truncated towards 0.
func NewFromFloat64(f float64, scale int) Decimal {
	i := big.NewInt(int64(f * math.Pow(float64(10), float64(scale))))
	return New(i, scale)
}

var ErrUnsupportedDecimalNotation = errors.New("luno: unsupported decimal notation")

func NewFromString(s string) (Decimal, error) {
	if strings.IndexByte(s, 'e') > -1 {
		return Decimal{}, ErrUnsupportedDecimalNotation
	}

	scale := getScale(s)
	s = strings.Replace(s, ".", "", 1)
	i, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return Decimal{}, ErrUnsupportedDecimalNotation
	}
	return New(i, scale), nil
}

// Zero returns a Decimal representing 0, with precision 0.
func Zero() Decimal {
	return New(cachedBigInt(0), 0)
}

// MarshalJSON converts the Decimal to JSON bytes.
func (d Decimal) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

// UnmarshalJSON reads JSON bytes into a Decimal.
func (d *Decimal) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	y, err := NewFromString(s)
	if err != nil {
		return err
	}
	*d = y
	return nil
}

// Strings returns a string representation of d.
func (d Decimal) String() string {
	di := bigIntDefault(d.i)
	s := di.String()

	var sign string
	if di.Sign() < 0 {
		sign = "-"
		s = strings.TrimPrefix(s, "-")
	}

	if d.scale > 0 {
		scale := int(d.scale)
		if len(s) < scale+1 {
			pad := scale + 1 - len(s)
			s = strings.Repeat("0", pad) + s
		}
		return fmt.Sprintf("%s%s.%s", sign, s[:len(s)-scale], s[len(s)-scale:])
	} else if d.scale < 0 {
		return sign + s + strings.Repeat("0", -int(d.scale))
	}

	return sign + s
}

func (d Decimal) QueryValue() string {
	return d.String()
}

// Float64 converts the decimal d to a float64.
func (d Decimal) Float64() float64 {
	if d.scale < 0 {
		d = d.ToScale(0)
	}

	scale := big.NewInt(int64(d.scale))
	scale.Exp(cachedBigInt(10), scale, nil)

	r := new(big.Rat)
	r.SetFrac(bigIntDefault(d.i), scale)

	f, _ := r.Float64()
	return f
}

// ToScale returns a Decimal representing the same value as d, but with the
// given scale. If scale is less than the Decimal's current scale, i.e. the new
// Decimal has fewer decimal points, the decimal is truncated (rounded towards
// 0). d is left unchanged.
func (d Decimal) ToScale(scale int) Decimal {
	exponent := d.scale - scale
	if exponent < 0 {
		exponent = -exponent
	}
	s := new(big.Int)
	s.Exp(cachedBigInt(10), cachedBigInt(int64(exponent)), nil)
	di := bigIntDefault(d.i)
	o := new(big.Int)
	if d.scale > scale {
		o.Quo(di, s)
	} else {
		o.Mul(di, s)
	}
	return Decimal{i: o, scale: scale}
}

// Sign returns -1 if d is negative, 1 if d is positive and 0 if d is zero.
func (d Decimal) Sign() int {
	return bigIntDefault(d.i).Sign()
}

// Cmp returns 1 if d>y, -1 if d<y or 0 if d==y.
func (d Decimal) Cmp(y Decimal) int {
	_d, _y := scaleToMax(d, y)
	di := bigIntDefault(_d.i)
	yi := bigIntDefault(_y.i)
	return di.Cmp(yi)
}

// Neg returns the negative of d. d is left unchanged.
func (d Decimal) Neg() Decimal {
	return Decimal{
		i:     new(big.Int).Neg(bigIntDefault(d.i)),
		scale: d.scale,
	}
}

// Add adds y to d and returns the result. d is left unchanged.
func (d Decimal) Add(y Decimal) Decimal {
	_d, _y := scaleToMax(d, y)
	di := bigIntDefault(_d.i)
	yi := bigIntDefault(_y.i)
	return New(new(big.Int).Add(di, yi), _d.scale)
}

// Sub subtracts y from d and returns the result. d is left unchanged.
func (d Decimal) Sub(y Decimal) Decimal {
	_d, _y := scaleToMax(d, y)
	di := bigIntDefault(_d.i)
	yi := bigIntDefault(_y.i)
	return New(new(big.Int).Sub(di, yi), _d.scale)
}

// MulInt64 multiplies d by y and returns the result. d is left unchanged.
func (d Decimal) MulInt64(y int64) Decimal {
	di := bigIntDefault(d.i)
	return New(new(big.Int).Mul(di, cachedBigInt(y)), d.scale)
}

// DivInt64 divides d by y and returns the result. d is left unchanged.
func (d Decimal) DivInt64(y int64) Decimal {
	di := bigIntDefault(d.i)
	return New(new(big.Int).Div(di, cachedBigInt(y)), d.scale)
}

// Mul multiplies d by y and returns the result. The result has a scale equal to
// the sum of d's and y's scales. d is left unchanged.
func (d Decimal) Mul(y Decimal) Decimal {
	di := bigIntDefault(d.i)
	yi := bigIntDefault(y.i)
	return New(new(big.Int).Mul(di, yi), d.scale+y.scale)
}

// Div divides d by y and returns the result in the provided scale. If the
// provided scale is too small for the result of d/y, the result is truncated
// towards 0. d is left unchanged.
func (d Decimal) Div(y Decimal, scale int) Decimal {
	if d.scale > y.scale+scale {
		y = y.ToScale(d.scale - scale)
	} else {
		d = d.ToScale(y.scale + scale)
	}
	return New(new(big.Int).Div(bigIntDefault(d.i), bigIntDefault(y.i)), scale)
}

func getScale(s string) int {
	i := strings.IndexByte(s, '.')
	if i == -1 {
		return 0
	}
	return len(s) - i - 1
}

func bigIntDefault(i *big.Int) *big.Int {
	if i == nil {
		return new(big.Int)
	}
	return i
}

func scaleToMax(x, y Decimal) (Decimal, Decimal) {
	if x.scale > y.scale {
		return x, y.ToScale(x.scale)
	} else if x.scale < y.scale {
		return x.ToScale(y.scale), y
	}
	return x, y
}
