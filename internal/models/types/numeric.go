package types

import "github.com/shopspring/decimal"

type Numeric string

const (
	NumericZero Numeric = Numeric("0")
)

func NewNumericFromStr(v string) (Numeric, error) {
	vv, err := decimal.NewFromString(v)
	if err != nil {
		return NumericZero, err
	}
	return Numeric(vv.String()), nil
}

func (a Numeric) String() string {
	return string(a)
}

func (a Numeric) Decimal() decimal.Decimal {
	vv, err := decimal.NewFromString(string(a))
	if err != nil {
		return decimal.Zero
	}
	return vv
}

// Cmp compares the numbers represented by d and d2 and returns:

// -1 if d <  d2
//
//	0 if d == d2
//
// +1 if d >  d2
func (d1 Numeric) Cmp(d2 Numeric) int {
	return d1.Decimal().Cmp(d2.Decimal())
}

func (d1 Numeric) Add(d2 Numeric) Numeric {
	return Numeric(d1.Decimal().Add(d2.Decimal()).String())
}

func (d1 Numeric) Sub(d2 Numeric) Numeric {
	return Numeric(d1.Decimal().Sub(d2.Decimal()).String())
}

func (d1 Numeric) Mul(d2 Numeric) Numeric {
	return Numeric(d1.Decimal().Mul(d2.Decimal()).String())
}

func (d1 Numeric) Div(d2 Numeric) Numeric {
	return Numeric(d1.Decimal().Div(d2.Decimal()).String())
}

func (d1 Numeric) Equal(d2 Numeric) bool {
	return d1.Decimal().Equal(d2.Decimal())
}

// Neg returns -d
func (d Numeric) Neg() Numeric {
	return Numeric(d.Decimal().Neg().String())
}
