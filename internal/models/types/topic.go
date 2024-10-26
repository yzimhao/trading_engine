package types

import (
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
)

const (
	TOPIC_ORDER_NEW            = "order_new" // %s 表示交易对symbol
	TOPIC_ORDER_CANCEL         = "order_cancel"
	TOPIC_PROCESS_ORDER_CANCEL = "process_order_cancel"
	TOPIC_ORDER_TRADE          = "order_trade"
	TOPIC_ORDER_SETTLE         = "order_settle"
)

type EventOrderNew struct {
	Symbol    string                   `json:"symbol"`
	OrderId   string                   `json:"order_id"`
	OrderSide matching_types.OrderSide `json:"order_side"`
	OrderType matching_types.OrderType `json:"order_type"`
	Price     *string                  `json:"price"`
	Quantity  *string                  `json:"quantity"`
	Amount    *string                  `json:"amount"`
	NanoTime  int64                    `json:"nano_time"`
}

type EventOrderCancel struct{}

type EventOrderTrade struct {
	matching_types.TradeResult
}

type EventOrderSettle struct {
	matching_types.TradeResult
}
