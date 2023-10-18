package www

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haobase/message"
	"github.com/yzimhao/trading_engine/cmd/haobase/message/ws"
	"github.com/yzimhao/trading_engine/types"
	"github.com/yzimhao/trading_engine/utils/app"
)

func publish_depth() {
	tsymbols := base.NewTSymbols()

	for _, item := range tsymbols.All() {
		push_depth_message(item.Symbol)
	}
}

func get_depth_data(symbol string) (*types.RedisDepthData, error) {
	topic := types.FormatDepthData.Format(symbol)

	rdc := app.RedisPool().Get()
	defer rdc.Close()

	raw, err := redis.Bytes(rdc.Do("GET", topic))
	if err != nil {
		return nil, err
	}

	var data types.RedisDepthData
	err = json.Unmarshal(raw, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func push_depth_message(symbol string) {
	go func() {
		for {
			func() (err error) {
				defer func() {
					if err != nil {
						logrus.Errorf("push_depth_message err: %s", err.Error())
						time.Sleep(time.Duration(3) * time.Second)
					}
				}()

				data, err := get_depth_data(symbol)
				//委托盘推送
				to_msg_depth := types.MsgDepth.Format(symbol)
				message.Publish(ws.MsgBody{
					To: to_msg_depth,
					Response: ws.Response{
						Type: to_msg_depth,
						Body: gin.H{
							"asks": data.Asks,
							"bids": data.Bids,
						},
					},
				})

				//最新价格
				to_latest_price := types.MsgLatestPrice.Format(symbol)
				message.Publish(ws.MsgBody{
					To: to_latest_price,
					Response: ws.Response{
						Type: to_latest_price,
						Body: map[string]any{
							"latest_price": data.Price,
							"at":           data.At,
						},
					},
				})

				//计算24H涨跌幅
				market_24h(symbol, data.Price)

				time.Sleep(time.Duration(1) * time.Second)
				return nil
			}()
		}
	}()
}
