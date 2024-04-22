package redisdb

import "encoding/json"

type OrderBookData struct {
	Price    string      `json:"price"`
	At       int64       `json:"at"`
	Asks     [][2]string `json:"asks"`
	Bids     [][2]string `json:"bids"`
	AsksSize int64       `json:"asks_size"`
	BidsSize int64       `json:"bids_size"`
}

func (r *OrderBookData) JSON() []byte {
	raw, _ := json.Marshal(r)
	return raw
}
