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
