package types

import (
	"github.com/shopspring/decimal"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
)

const (
	TOPIC_ORDER_NEW            = "order_new" // %s 表示交易对symbol
	TOPIC_NOTIFY_ORDER_CANCEL  = "notify_order_cancel"
	TOPIC_PROCESS_ORDER_CANCEL = "process_order_cancel"
	TOPIC_ORDER_TRADE          = "order_trade"
	TOPIC_ORDER_SETTLE         = "order_settle"
	TOPIC_NOTIFY_QUOTE         = "notify_quote"
)

type EventOrderNew struct {
	Symbol    string                   `json:"symbol,omitempty"`
	OrderId   string                   `json:"order_id,omitempty"`
	OrderSide matching_types.OrderSide `json:"order_side,omitempty"`
	OrderType matching_types.OrderType `json:"order_type,omitempty"`
	Price     *decimal.Decimal         `json:"price,omitempty"`
	Quantity  *decimal.Decimal         `json:"quantity,omitempty"`
	Amount    *decimal.Decimal         `json:"amount,omitempty"`
	MaxAmount *decimal.Decimal         `json:"max_amount,omitempty"`
	MaxQty    *decimal.Decimal         `json:"max_qty,omitempty"`
	NanoTime  int64                    `json:"nano_time,omitempty"`
}

type EventNotifyCancelOrder struct {
	Symbol    string                    `json:"symbol,omitempty"`
	OrderSide matching_types.OrderSide  `json:"order_side,omitempty"`
	OrderId   string                    `json:"order_id,omitempty"`
	Type      matching_types.RemoveType `json:"type,omitempty"`
}

type EventOrderTrade struct {
	matching_types.TradeResult
}

type EventOrderSettle struct {
	matching_types.TradeResult
}

type EventNotifyQuote struct {
	matching_types.TradeResult
	// Symbol        string                 `json:"symbol,omitempty"`
	// AskOrderId    string                 `json:"ask,omitempty"`
	// BidOrderId    string                 `json:"bid,omitempty"`
	// TradeQuantity decimal.Decimal        `json:"trade_quantity,omitempty"`
	// TradePrice    decimal.Decimal        `json:"trade_price,omitempty"`
	// TradeBy       matching_types.TradeBy `json:"trade_by,omitempty"`
	// TradeTime     int64                  `json:"trade_time,omitempty"`
}

type EventCancelOrder struct {
	Symbol  string `json:"symbol,omitempty"`
	OrderId string `json:"order_id,omitempty"`
}
