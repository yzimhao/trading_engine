package settle

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gookit/goutil/arrutil"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haobase/orders"
	"github.com/yzimhao/trading_engine/config"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/types/redisdb"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
	"github.com/yzimhao/trading_engine/utils/app/keepalive"
)

func Run() {
	for {
		init_symbols()
		time.Sleep(time.Second * 5)
	}
}

func init_symbols() {
	local_config_symbols := config.App.Local.Symbols
	db_symbols := base.NewTradeSymbol().All()
	for _, item := range db_symbols {
		if len(local_config_symbols) > 0 && arrutil.Contains(local_config_symbols, item.Symbol) || len(local_config_symbols) == 0 {
			if !keepalive.HasExtrasKeyValue("settle.symbols", item.Symbol) {
				run_clearing(item.Symbol)
				keepalive.SetExtras("settle.symbols", item.Symbol)
			}
		}
	}
}

func run_clearing(symbol string) {
	//成交日志队列
	go watch_tradeok_list(symbol)
}

func watch_tradeok_list(symbol string) {
	key := redisdb.TradeResultQueue.Format(redisdb.Replace{"symbol": symbol})
	app.Logger.Infof("监听%s成交日志，等待结算...", symbol)
	for {
		func() {
			rdc := app.RedisPool().Get()
			defer rdc.Close()

			if n, _ := redis.Int64(rdc.Do("LLen", key)); n == 0 {
				time.Sleep(time.Duration(50) * time.Millisecond)
				return
			}

			raw, _ := redis.Bytes(rdc.Do("Lpop", key))
			app.Logger.Infof("收到%s成交记录: %s", symbol, raw)
			go clearing_trade_order(symbol, raw)
		}()

	}
}

func clearing_trade_order(symbol string, raw []byte) {
	var data trading_core.TradeResult
	err := json.Unmarshal(raw, &data)
	if err != nil {
		app.Logger.Errorf("%s成交日志格式错误: %s %s", symbol, err.Error(), raw)
		return
	}

	orders.Lock(orders.SettleLock, data.AskOrderId)
	orders.Lock(orders.SettleLock, data.BidOrderId)

	if data.Last == "" {
		go newClean(data)
	} else {
		go func() {
			for {
				time.Sleep(time.Duration(50) * time.Millisecond)
				app.Logger.Infof("等待订单%s 其他成交结算完成...", data.Last)
				if orders.GetLock(orders.SettleLock, data.Last) == 1 {
					newClean(data)
					break
				}
			}
		}()
	}

}

func generate_trading_id(ask, bid string) string {
	times := time.Now().Format("060102")
	hash := utils.Hash256(fmt.Sprintf("%s%s", ask, bid))
	return fmt.Sprintf("T%s%s", times, hash[0:17])
}

func notify_quote(raw trading_core.TradeResult) {
	//通知kline系统
	rdc := app.RedisPool().Get()
	defer rdc.Close()

	quote_key := redisdb.QuoteTradeResultQueue.Format(redisdb.Replace{"symbol": raw.Symbol})
	if _, err := rdc.Do("RPUSH", quote_key, raw.Json()); err != nil {
		app.Logger.Errorf("RPUSH %s err: %s", quote_key, err.Error())
	}
}
