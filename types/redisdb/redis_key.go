package redisdb

import (
	"strings"

	"github.com/yzimhao/trading_engine/config"
)

type redisdb string

const (
	Keepalive             redisdb = "{prefix}keepalive.{uuid}"
	NewOrderQueue         redisdb = "{prefix}queue.new.order.{symbol}"
	CancelOrderQueue      redisdb = "{prefix}queue.cancel.order.{symbol}"
	TradeResultQueue      redisdb = "{prefix}queue.trade.result.{symbol}"
	CancelResultQueue     redisdb = "{prefix}queue.cancel.result.{symbol}"
	QuoteTradeResultQueue redisdb = "{prefix}queue.quote.trade.result.{symbol}"
	WsMessageQueue        redisdb = "{prefix}queue.ws.message"
	//#
	OrderBook redisdb = "{prefix}orderbook.{symbol}"
	//订单结算
	OrderLock redisdb = "{prefix}order.lock.{order_id}"

	BroadcastLatestPrice redisdb = "{prefix}broadcast.latest_price.{symbol}"

	OrderDetail            redisdb = "{prefix}order.detail.{order_id}"
	SymbolUnfinishedOrders redisdb = "{prefix}unfinished.order.{symbol}"
	SymbolLatestPrice      redisdb = "{prefix}latest_price.{symbol}"

	//base
	BaseTradeSymbolAll redisdb = "{prefix}base.trade.symbol.all"
)

type Replace map[string]string

func (r redisdb) Format(kv Replace) string {
	str := strings.Replace(string(r), "{prefix}", config.App.Redis.Prefix, -1)

	for k, v := range kv {
		str = strings.Replace(str, "{"+k+"}", v, -1)
	}

	return str
}
