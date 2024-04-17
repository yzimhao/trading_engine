package redisdb

import "encoding/json"

type OrderBookData struct {
	Price string      `json:"price"`
	At    int64       `json:"at"`
	Asks  [][2]string `json:"asks"`
	Bids  [][2]string `json:"bids"`
}

func (r *OrderBookData) JSON() []byte {
	raw, _ := json.Marshal(r)
	return raw
}
