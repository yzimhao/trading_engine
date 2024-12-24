package webws

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

var (
	M *Hub
)

type subMessage struct {
	Subsc   []string `json:"sub"`
	UnSubsc []string `json:"unsub"`
}

type MsgBody struct {
	To       string
	Response Response
}

type Response struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

func (m *MsgBody) BodyHash() string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%v", m.Response)))
	return hex.EncodeToString(h.Sum(nil))
}

func (m *MsgBody) GetBody() []byte {
	re := m.Response
	data, _ := json.Marshal(re)
	return data
}

func (m *MsgBody) JSON() []byte {
	raw, _ := json.Marshal(m)
	return raw
}

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
