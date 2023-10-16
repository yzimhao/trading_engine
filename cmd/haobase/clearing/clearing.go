package clearing

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haobase/base/symbols"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/types"
	"github.com/yzimhao/trading_engine/utils"
)

func Run() {
	//load symbols
	db := base.DB().NewSession()
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
			cx := context.Background()
			rdc := base.RDC()
			defer rdc.Close()

			if n, _ := rdc.LLen(cx, key).Result(); n == 0 {
				time.Sleep(time.Duration(50) * time.Millisecond)
				return
			}

			raw, _ := rdc.LPop(cx, key).Bytes()

			var data trading_core.TradeResult
			err := json.Unmarshal(raw, &data)
			if err != nil {
				logrus.Warnf("%s 解析json: %s 错误: %s", key, raw, err)
				return
			}

			logrus.Infof("%s成交记录 ask: %s bid: %s price: %s vol: %s", data.Symbol, data.AskOrderId, data.BidOrderId, data.TradePrice.String(), data.TradeQuantity.String())

			err = newClean(data)
			if err != nil {
				logrus.Warnf("结算错误: %s %s", raw, err.Error())
				return
			}

			//通知kline系统
			rdc.RPush(cx, quote_key, raw)

			// if !data.Last {
			// 	go newClean(data)
			// } else {

			// }
		}()

	}
}

func generate_trading_id(ask, bid string) string {
	times := time.Now().Format("060102")
	hash := utils.Hash256(fmt.Sprintf("%s%s", ask, bid))
	return fmt.Sprintf("T%s%s", times, hash[0:17])
}
