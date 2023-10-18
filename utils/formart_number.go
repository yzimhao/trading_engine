package utils

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type FloatString string

func (t FloatString) MarshalJSON() ([]byte, error) {
	tt, _ := decimal.NewFromString(string(t))
	s := fmt.Sprintf("\"%s\"", tt.String())
	return []byte(s), nil
}

func (t FloatString) String() string {
	return string(t)
}

func (t FloatString) Decimal() decimal.Decimal {
	tt, _ := decimal.NewFromString(string(t))
	return tt
}
