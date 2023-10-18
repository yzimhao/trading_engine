package types

import (
	"strings"

	"github.com/spf13/viper"
)

type RedisKey string

const (
	FormatNewOrder             RedisKey = "{prefix}order.list.{symbol}"
	FormatCancelOrder          RedisKey = "{prefix}need.cancel.{symbol}"
	FormatTradeResult          RedisKey = "{prefix}trade.result.{symbol}"
	FormatCancelResult         RedisKey = "{prefix}cancel.result.{symbol}"
	FormatQuoteTradeResult     RedisKey = "{prefix}quote.trade.result.{symbol}"
	FormatWsMessage            RedisKey = "{prefix}ws.message"
	FormatBroadcastDepth       RedisKey = "{prefix}broadcast.depth.{symbol}"
	FormatBroadcastLatestPrice RedisKey = "{prefix}broadcast.latest_price.{symbol}"
)

func (r RedisKey) Format(symbol string) string {
	key := strings.Replace(r.String(), "{symbol}", symbol, -1)
	key = strings.Replace(key, "{prefix}", viper.GetString("redis.prefix"), -1)
	return key
}

func (r RedisKey) String() string {
	return string(r)
}

type ChannelLatestPrice struct {
	T     int64  `json:"t"`
	Price string `json:"price"`
}
