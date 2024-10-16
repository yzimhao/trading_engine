package types

import "github.com/shopspring/decimal"

type Amount string

func (a Amount) NewFromStr(v string) (Amount, error) {
	vv, err := decimal.NewFromString(v)
	if err != nil {
		return "0", err
	}
	return Amount(vv.String()), nil
}

func (a Amount) String() string {
	return string(a)
}

// Cmp compares the numbers represented by d and d2 and returns:

// -1 if d <  d2
//
//	0 if d == d2
//
// +1 if d >  d2
func (d1 Amount) Cmp(d2 Amount) int {
	aa, _ := decimal.NewFromString(string(d1))
	bb, _ := decimal.NewFromString(string(d2))
	return aa.Cmp(bb)
}

func (d1 Amount) Add(d2 Amount) Amount {
	aa, _ := decimal.NewFromString(string(d1))
	bb, _ := decimal.NewFromString(string(d2))
	return Amount(aa.Add(bb).String())
}

func (d1 Amount) Sub(d2 Amount) Amount {
	aa, _ := decimal.NewFromString(string(d1))
	bb, _ := decimal.NewFromString(string(d2))
	return Amount(aa.Sub(bb).String())
}

func (d1 Amount) Mul(d2 Amount) Amount {
	aa, _ := decimal.NewFromString(string(d1))
	bb, _ := decimal.NewFromString(string(d2))
	return Amount(aa.Mul(bb).String())
}

func (d1 Amount) Div(d2 Amount) Amount {
	aa, _ := decimal.NewFromString(string(d1))
	bb, _ := decimal.NewFromString(string(d2))
	if bb.IsZero() {
		//TODO 这里需要抛出错误
		return "0"
	}
	return Amount(aa.Div(bb).String())
}
