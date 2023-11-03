package utils

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type DeciStr string

func (t DeciStr) MarshalJSON() ([]byte, error) {
	tt, _ := decimal.NewFromString(string(t))
	s := fmt.Sprintf("\"%s\"", tt.String())
	return []byte(s), nil
}

func (t DeciStr) String() string {
	return string(t)
}

func (t DeciStr) Decimal() decimal.Decimal {
	tt, _ := decimal.NewFromString(string(t))
	return tt
}

func FormatDecimal(d string, digit int) string {
	// f, _ := D(d).Float64()
	// format := "%." + fmt.Sprintf("%d", digit) + "f"
	// return fmt.Sprintf(format, f)

	// return D(d).Truncate(int32(digit)).String()
	return D(d).StringFixedBank(int32(digit))
}
