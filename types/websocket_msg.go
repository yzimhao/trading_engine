package types

import "strings"

type WebSocketMsgType string

const (
	MsgDepth       WebSocketMsgType = "depth.{symbol}"
	MsgTrade       WebSocketMsgType = "trade.{symbol}"
	MsgLatestPrice WebSocketMsgType = "price.{symbol}"
	MsgMarketKLine WebSocketMsgType = "kline.{period}.{symbol}"
	MsgMarket24H   WebSocketMsgType = "market.24h.{symbol}"
	MsgOrderCancel WebSocketMsgType = "order.cancel.{symbol}"
	MsgToken       WebSocketMsgType = "token.{token}"
	MsgUser        WebSocketMsgType = "_user.{user_id}" //特殊的类型，通过后端程序设置的属性
)

var (
	AllWebSocketMsg = []WebSocketMsgType{
		MsgDepth, MsgDepth, MsgLatestPrice, MsgMarketKLine, MsgMarket24H,
	}
)

func (w WebSocketMsgType) Format(data map[string]string) string {
	nw := string(w)
	for k, v := range data {
		nw = strings.Replace(nw, "{"+k+"}", v, -1)
	}
	return nw
}
