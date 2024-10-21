package types

import (
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
)

const (
	TOPIC_ORDER_NEW    = "order_new" // %s 表示交易对symbol
	TOPIC_ORDER_CANCEL = "order_cancel"
	TOPIC_ORDER_TRADE  = "order_trade"
	TOPIC_ORDER_SETTLE = "order_settle"
)

type EventOrderNew struct {
	Symbol string `json:"symbol"`
	At     int64  `json:"at"`
	// ...
}

type EventOrderCancel struct{}

type EventOrderTrade struct {
	matching_types.TradeResult
}

type EventOrderSettle struct {
	matching_types.TradeResult
}
