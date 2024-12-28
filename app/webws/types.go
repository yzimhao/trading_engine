package webws

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

type RecviceTag struct {
	Subscribe   []string `json:"subscribe,omitempty"`
	Unsubscribe []string `json:"unsubscribe,omitempty"`
}

type Response struct {
	Type string `json:"type,omitempty"`
	Body any    `json:"body,omitempty"`
}

type Message struct {
	To       string   `json:"to"`
	Response Response `json:"response"`
}

func (m *Message) Sign() string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%v", m.Response.Body)))
	return hex.EncodeToString(h.Sum(nil))
}

func (m *Message) ResponseBytes() []byte {
	raw, _ := json.Marshal(m.Response)
	return raw
}

func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func NewMessage(to string, tag string, body any) Message {
	m := Message{
		To: to,
		Response: Response{
			Type: tag,
			Body: body,
		},
	}
	return m
}

// websocket message tags
type MessageTagTpl string

const (
	MsgDepthTpl       MessageTagTpl = "depth.{symbol}"
	MsgTradeTpl       MessageTagTpl = "trade.{symbol}"
	MsgLatestPriceTpl MessageTagTpl = "price.{symbol}"
	MsgMarketKLineTpl MessageTagTpl = "kline.{period}.{symbol}"
	MsgMarket24HTpl   MessageTagTpl = "market.24h.{symbol}"
	MsgOrderCancelTpl MessageTagTpl = "order.cancel.{symbol}"
	MsgNewOrderTpl    MessageTagTpl = "order.new.{symbol}"
	MsgTokenTpl       MessageTagTpl = "token.{token}"
	MsgUserTpl        MessageTagTpl = "_user.{user_id}" //特殊的类型，通过后端程序设置的属性
)

var (
	AllWebSocketMsg = []MessageTagTpl{
		MsgDepthTpl, MsgTradeTpl, MsgLatestPriceTpl,
		MsgMarketKLineTpl, MsgMarket24HTpl,
		MsgOrderCancelTpl, MsgNewOrderTpl,
	}
)

func (w MessageTagTpl) Format(data map[string]string) string {
	nw := string(w)
	for k, v := range data {
		nw = strings.Replace(nw, "{"+k+"}", v, -1)
	}
	return nw
}
