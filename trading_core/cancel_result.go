package trading_core

import "encoding/json"

type CancelType int

const (
	CancelTypeByUnknown CancelType = iota + 1 //未知原因
	CancelTypeByUser                          //用户取消
	CancelTypeBySystem                        //系统自动取消
	CancelTypeByPartial                       //部分成交，剩余取消
	CancelTypeByMarket                        //市场条件变化 导致取消
)

type CancelBody struct {
	OrderId string     `json:"order_id"`
	Reason  CancelType `json:"reason"`
}

func (t *CancelBody) Json() []byte {
	raw, _ := json.Marshal(&t)
	return raw
}
