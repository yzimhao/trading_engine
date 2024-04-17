package www

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haobase/message"
	"github.com/yzimhao/trading_engine/cmd/haobase/message/ws"
	"github.com/yzimhao/trading_engine/types"
	"github.com/yzimhao/trading_engine/types/redisdb"
	"github.com/yzimhao/trading_engine/utils/app"
)

func publish_depth() {
	tsymbols := base.NewTradeSymbol()

	for _, item := range tsymbols.All() {
		push_depth_message(item.Symbol)
	}
}

func get_depth_data(symbol string) (*redisdb.OrderBookData, error) {
	topic := redisdb.DepthData.Format(redisdb.Replace{"symbol": symbol})

	rdc := app.RedisPool().Get()
	defer rdc.Close()

	raw, err := redis.Bytes(rdc.Do("GET", topic))
	if err != nil {
		return nil, err
	}

	var data redisdb.OrderBookData
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
						app.Logger.Errorf("push_depth_message err: %s", err.Error())
						time.Sleep(time.Duration(3) * time.Second)
					}
				}()

				data, err := get_depth_data(symbol)
				if err != nil {
					return
				}
				//委托盘推送
				to_msg_depth := types.MsgDepth.Format(map[string]string{
					"symbol": symbol,
				})
				message.Publish(ws.MsgBody{
					To: to_msg_depth,
					Response: ws.Response{
						Type: to_msg_depth,
						Body: gin.H{
							"asks": limitSize(data.Asks, 10),
							"bids": limitSize(data.Bids, 10),
						},
					},
				})

				//最新价格
				to_latest_price := types.MsgLatestPrice.Format(map[string]string{
					"symbol": symbol,
				})
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

				time.Sleep(time.Duration(300) * time.Millisecond)
				return nil
			}()
		}
	}()
}
