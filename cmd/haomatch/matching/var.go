package matching

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

type Order struct {
	OrderId   string `json:"order_id"`
	OrderType string `json:"order_type"`
	Side      string `json:"side"`
	Price     string `json:"price,omitempty"`
	Qty       string `json:"qty,omitempty"`
	Amount    string `json:"amount,omitempty"`
	MaxQty    string `json:"max_qty,omitempty"`
	MaxAmount string `json:"max_amount,omitempty"`
	At        int64  `json:"at"`
}

type cancel_order struct {
	Side    string `json:"side"`
	OrderId string `json:"order_id"`
}

func (o *Order) Json() []byte {
	raw, _ := json.Marshal(o)
	return raw
}

func d(a string) decimal.Decimal {
	d, _ := decimal.NewFromString(a)
	return d
}
