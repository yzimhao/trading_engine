package matching

import (
	"encoding/json"

	"github.com/yzimhao/trading_engine/trading_core"
)

type Order struct {
	OrderId   string                 `json:"order_id"`
	OrderType trading_core.OrderType `json:"order_type"`
	Side      trading_core.OrderSide `json:"side"`
	Price     string                 `json:"price,omitempty"`
	Qty       string                 `json:"qty,omitempty"`
	Amount    string                 `json:"amount,omitempty"`
	MaxQty    string                 `json:"max_qty,omitempty"`
	MaxAmount string                 `json:"max_amount,omitempty"`
	At        int64                  `json:"at"`
}

func (o *Order) Json() []byte {
	raw, _ := json.Marshal(o)
	return raw
}

type StructCancelOrder struct {
	Side    trading_core.OrderSide  `json:"side"`
	OrderId string                  `json:"order_id"`
	Reason  trading_core.CancelType `json:"reason"`
	Status  string                  `json:"status"`
}

func (o *StructCancelOrder) Json() []byte {
	raw, _ := json.Marshal(o)
	return raw
}
