package clearing

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/yzimhao/trading_engine/cmd/haobase/base/symbols"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/types"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func Run() {
	//load symbols
	db := app.Database().NewSession()
	defer db.Close()

	var rows []symbols.TradingVarieties
	db.Table(new(symbols.TradingVarieties)).Find(&rows)

	for _, row := range rows {
		run_clearing(row.Symbol)
	}
}

func run_clearing(symbol string) {
	go watch_redis_list(symbol)
}

func watch_redis_list(symbol string) {
	key := types.FormatTradeResult.Format(symbol)
	quote_key := types.FormatQuoteTradeResult.Format(symbol)
	logrus.Infof("结算，正在监听%s成交日志...", symbol)
	for {
		func() {
			rdc := app.RedisPool().Get()
			defer rdc.Close()

			if n, _ := redis.Int64(rdc.Do("LLen", key)); n == 0 {
				time.Sleep(time.Duration(50) * time.Millisecond)
				return
			}

			raw, _ := redis.Bytes(rdc.Do("Lpop", key))

			var data trading_core.TradeResult
			err := json.Unmarshal(raw, &data)
			if err != nil {
				logrus.Warnf("%s 解析json: %s 错误: %s", key, raw, err)
				return
			}

			logrus.Infof("%s成交记录 ask: %s bid: %s price: %s vol: %s", data.Symbol, data.AskOrderId, data.BidOrderId, data.TradePrice.String(), data.TradeQuantity.String())

			if data.Last {
				//todo 优化
				time.Sleep(time.Duration(50) * time.Millisecond)
				err = newClean(data)
			} else {
				err = newClean(data)
			}
			if err != nil {
				logrus.Warnf("结算错误: %s %s", raw, err.Error())
				return
			}

			//通知kline系统
			if _, err := rdc.Do("RPUSH", quote_key, raw); err != nil {
				logrus.Errorf("rpush %s err: %s", quote_key, err.Error())
			}
		}()

	}
}

func generate_trading_id(ask, bid string) string {
	times := time.Now().Format("060102")
	hash := utils.Hash256(fmt.Sprintf("%s%s", ask, bid))
	return fmt.Sprintf("T%s%s", times, hash[0:17])
}
