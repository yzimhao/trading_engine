package trading_core

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

type TradeBy int

const (
	BySeller TradeBy = 1
	ByBuyer  TradeBy = 2
)

type TradeResult struct {
	Symbol        string          `json:"symbol"`
	AskOrderId    string          `json:"ask"`
	BidOrderId    string          `json:"bid"`
	TradeQuantity decimal.Decimal `json:"trade_quantity"`
	TradePrice    decimal.Decimal `json:"trade_price"`
	TradeBy       TradeBy         `json:"trade_by"`
	TradeTime     int64           `json:"trade_time"` //纳秒级
	Last          string          `json:"last"`       //市价订单标记，用于结算时取消市价单剩余未成交的部分
}

func (t *TradeResult) Json() []byte {
	raw, _ := json.Marshal(&t)
	return raw
}
