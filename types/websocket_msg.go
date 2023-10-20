package types

import "strings"

type WebSocketMsgType string
type WebSocketKlineMsgType string
type WebSocketMarket string

const (
	MsgDepth           WebSocketMsgType      = "depth.{symbol}"
	MsgTrade           WebSocketMsgType      = "tradelog.{symbol}"
	MsgLatestPrice     WebSocketMsgType      = "latest_price.{symbol}"
	MsgMarketKLine     WebSocketKlineMsgType = "kline.{period}.{symbol}"
	MsgMarket24H       WebSocketMarket       = "market.24h.{symbol}"
	MsgUserOrderCancel WebSocketMsgType      = "user.order.cancel.{symbol}"
)

func (w WebSocketMsgType) Format(symbol string) string {
	key := strings.Replace(string(w), "{symbol}", symbol, -1)
	return key
}

func (w WebSocketKlineMsgType) Format(period, symbol string) string {
	key := strings.Replace(string(w), "{symbol}", symbol, -1)
	key = strings.Replace(key, "{period}", period, -1)
	return key
}

func (w WebSocketMarket) Format(symbol string) string {
	key := strings.Replace(string(w), "{symbol}", symbol, -1)
	return key
}
