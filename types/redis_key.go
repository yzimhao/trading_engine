package types

import (
	"encoding/json"
	"strings"

	"github.com/yzimhao/trading_engine/config"
)

type RedisKey string

const (
	FormatNewOrder             RedisKey = "{prefix}order.list.{symbol}"
	FormatCancelOrder          RedisKey = "{prefix}need.cancel.{symbol}"
	FormatTradeResult          RedisKey = "{prefix}trade.result.{symbol}"
	FormatCancelResult         RedisKey = "{prefix}cancel.result.{symbol}"
	FormatQuoteTradeResult     RedisKey = "{prefix}quote.trade.result.{symbol}"
	FormatWsMessage            RedisKey = "{prefix}ws.message"
	FormatDepthData            RedisKey = "{prefix}depth.{symbol}"
	FormatBroadcastLatestPrice RedisKey = "{prefix}broadcast.latest_price.{symbol}"
)

func (r RedisKey) Format(symbol string) string {
	key := strings.Replace(r.String(), "{symbol}", symbol, -1)
	key = strings.Replace(key, "{prefix}", config.App.Redis.Prefix, -1)
	return key
}

func (r RedisKey) String() string {
	return string(r)
}

type RedisDepthData struct {
	Price string      `json:"price"`
	At    int64       `json:"at"`
	Asks  [][2]string `json:"asks"`
	Bids  [][2]string `json:"bids"`
}

func (r *RedisDepthData) JSON() []byte {
	raw, _ := json.Marshal(r)
	return raw
}
