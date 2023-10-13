package clearing

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/types"
	"github.com/yzimhao/trading_engine/utils"
)

func RunClearing(symbol string) {
	watch_redis_list(symbol)
}

func watch_redis_list(symbol string) {
	key := types.FormatTradeResult.Format(symbol)
	logrus.Infof("结算，正在监听%s成交日志...", symbol)
	for {
		func() {
			cx := context.Background()
			rdc := base.RDC()

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

			newClean(data)
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
