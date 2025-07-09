package types

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

type OrderType string
type OrderSide string
type SubOrderType string
type TradeBy string

const (
	OrderTypeLimit             OrderType    = "LIMIT"
	OrderTypeMarket            OrderType    = "MARKET"
	SubOrderTypeMarketByQty    SubOrderType = "marketByQty"
	SubOrderTypeMarketByAmount SubOrderType = "marketByAmount"
	SubOrderTypeUnknown        SubOrderType = "unknown"
	OrderSideBuy               OrderSide    = "BUY"
	OrderSideSell              OrderSide    = "SELL"
	TradeBySeller              TradeBy      = "SELLER"
	TradeByBuyer               TradeBy      = "BUYER"
)

func (os OrderSide) String() string {
	return string(os)
}

func (ot OrderType) String() string {
	return string(ot)
}

type TradeResult struct {
	Symbol          string           `json:"symbol"`
	AskOrderId      string           `json:"ask"`
	BidOrderId      string           `json:"bid"`
	TradeQuantity   decimal.Decimal  `json:"trade_quantity"`
	TradePrice      decimal.Decimal  `json:"trade_price"`
	TradeBy         TradeBy          `json:"trade_by"`
	TradeTime       int64            `json:"trade_time"`        //纳秒级
	MarketOrderInfo *MarketOrderInfo `json:"market_order_info"` //市价订单标记，用于结算时取消市价单剩余未成交的部分
}

type MarketOrderInfo struct {
	OrderId      string          `json:"order_id"`            // 市价订单ID
	IsFinalTrade bool            `json:"is_final_trade"`      // 是否为该市价订单的最后一笔成交
	Remaining    decimal.Decimal `json:"remaining,omitempty"` // 剩余未成交量(可选)
}

func (t *TradeResult) MarshalBinary() (data []byte, err error) {
	data, err = json.Marshal(&t)
	return
}

func (t *TradeResult) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}

type RemoveResult struct {
	Symbol   string         `json:"symbol"`
	UniqueId string         `json:"unique_id"`
	Type     RemoveItemType `json:"type"`
}
