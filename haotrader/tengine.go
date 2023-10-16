package haotrader

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/types"
	"github.com/yzimhao/trading_engine/utils/app"
)

type tengine struct {
	symbol              string
	tp                  *trading_core.TradePair
	restore_done_signal chan struct{}
	broadcast           chan struct{}
}

func NewTengine(symbol string, price_digit, qty_digit int) *trading_core.TradePair {
	te := tengine{
		symbol:              symbol,
		tp:                  trading_core.NewTradePair(symbol, price_digit, qty_digit),
		restore_done_signal: make(chan struct{}),
		broadcast:           make(chan struct{}, 100),
	}

	go te.queue_monitor()
	go te.restore()
	go te.broadcast_depth()
	go te.pull_new_order()
	go te.pull_cancel_order()
	go te.monitor_result()
	return te.tp
}

func (t *tengine) broadcast_depth() {
	depth_channel := types.FormatBroadcastDepth.Format(t.symbol)
	price_channel := types.FormatBroadcastLatestPrice.Format(t.symbol)
	//如果长时间没有触发，5s自动触发一次更新
	go func() {
		for {
			time.Sleep(time.Duration(5) * time.Second)
			t.broadcast <- struct{}{}
		}
	}()

	for {
		select {
		case <-t.broadcast:
			go func() {
				data := gin.H{
					//todo 限制最大获取的数量
					"asks": t.tp.GetAskDepth(100),
					"bids": t.tp.GetBidDepth(100),
				}

				rdc := app.RedisPool().Get()
				defer rdc.Close()

				raw, _ := json.Marshal(data)

				if _, err := rdc.Do("Publish", depth_channel, raw); err != nil {
					logrus.Warnf("广播%s消息失败: %s", depth_channel, err)
				}
			}()

			go func() {
				price, at := t.tp.LatestPrice()
				data := types.ChannelLatestPrice{
					T:     at,
					Price: t.tp.Price2String(price),
				}

				raw, _ := json.Marshal(data)
				rdc := app.RedisPool().Get()
				defer rdc.Close()
				if _, err := rdc.Do("Publish", price_channel, raw); err != nil {
					logrus.Warnf("广播%s消息失败: %s", price_channel, err)
				}
			}()

		default:
			time.Sleep(time.Duration(100) * time.Millisecond)
		}
	}
}

func (t *tengine) queue_monitor() {
	t.tp.OnEvent(func(qi trading_core.QueueItem) {
		//恢复数据完成后，再开始数据持久化
		if t.tp.TriggerEvent() {
			raw := Order{
				OrderId:   qi.GetUniqueId(),
				Side:      qi.GetOrderSide().String(),
				OrderType: "limit",
				Price:     qi.GetPrice().String(),
				Qty:       qi.GetQuantity().String(),
				At:        qi.GetCreateTime(),
			}

			if qi.GetQuantity().Cmp(decimal.Zero) > 0 {
				logrus.Debugf("queue event update: %#v", raw)
				go localdb.Set(t.symbol, raw.OrderId, raw.Json())
			}
		}
		t.broadcast <- struct{}{}
	}, func(qi trading_core.QueueItem) {
		if t.tp.TriggerEvent() {
			raw := Order{
				OrderId:   qi.GetUniqueId(),
				Side:      qi.GetOrderSide().String(),
				OrderType: "limit",
				Price:     qi.GetPrice().String(),
				Qty:       "0",
				At:        qi.GetCreateTime(),
			}
			logrus.Debugf("queue event remove: %#v", raw)
			go localdb.Remove(t.symbol, raw.OrderId)
		}

		t.broadcast <- struct{}{}
	}, func(tl trading_core.TradeResult) {
		//只保留最近的1条成交记录,用于恢复最新成交价格
		localdb.Set("tradelog", t.symbol, tl.Json())
	})
}

func (t *tengine) restore() {

	defer func() {
		logrus.Infof("[%s]数据恢复 已完成", t.symbol)
		close(t.restore_done_signal)
		t.tp.SetTriggerEvent(true)
		t.tp.SetPauseMatch(false)
	}()
	//从磁盘恢复上一次的数据，先暂停撮合系统的撮合，等数据全部加载完成后再开启撮合
	t.tp.SetPauseMatch(true)

	//恢复orderbook
	data := localdb.Find(t.symbol, "")
	logrus.Infof("正在恢复[%s]数据，共%d条", t.symbol, len(data))
	for i, v := range data {
		func(n int, raw []byte) {
			logrus.Infof("恢复数据[%s]第%d条: %s", t.symbol, n+1, raw)
			var data Order
			json.Unmarshal(raw, &data)

			if data.Side == "ask" {
				t.tp.PushNewOrder(trading_core.NewAskLimitItem(data.OrderId, d(data.Price), d(data.Qty), data.At))
			} else {
				t.tp.PushNewOrder(trading_core.NewBidLimitItem(data.OrderId, d(data.Price), d(data.Qty), data.At))
			}
		}(i, v)
	}

	//恢复最新成交价格
	tls := localdb.Find("tradelog", t.symbol)
	for _, v := range tls {
		var tlog trading_core.TradeResult
		json.Unmarshal(v, &tlog)
		t.tp.SetLatestPrice(tlog.TradePrice)
	}

}

func (t *tengine) pull_new_order() {
	<-t.restore_done_signal
	key := types.FormatNewOrder.Format(t.symbol)
	logrus.Infof("正在监听redis队列: %s", key)
	for {

		// cx := context.Background()
		// if n, _ := rdc.LLen(cx, key).Result(); n == 0 || t.tp.IsPausePushNew() {
		// 	time.Sleep(time.Duration(50) * time.Millisecond)
		// 	continue
		// }
		func() {

			rdc := app.RedisPool().Get()
			defer rdc.Close()
			if n, _ := redis.Int64(rdc.Do("LLen", key)); n == 0 || t.tp.IsPausePushNew() {
				time.Sleep(time.Duration(50) * time.Millisecond)
				return
			}

			raw, _ := redis.Bytes(rdc.Do("Lpop", key))

			go func(raw []byte) {
				var data Order
				err := json.Unmarshal(raw, &data)
				if err != nil {
					logrus.Warnf("%s 解析json: %s 错误: %s", key, raw, err)
				}

				if data.OrderId != "" {
					logrus.Debugf("%s队列LPop: %s", key, raw)
					side := strings.ToLower(data.Side)
					order_type := strings.ToLower(data.OrderType)

					if order_type == "limit" {
						if side == trading_core.OrderSideSell.String() {
							t.tp.PushNewOrder(trading_core.NewAskLimitItem(data.OrderId, d(data.Price), d(data.Qty), data.At))
						} else if side == trading_core.OrderSideBuy.String() {
							t.tp.PushNewOrder(trading_core.NewBidLimitItem(data.OrderId, d(data.Price), d(data.Qty), data.At))
						} else {
							logrus.Errorf("新订单参数错误: %s side只能是sell/buy", raw)
						}
					} else if order_type == "market_qty" {
						//按成交量
						if side == trading_core.OrderSideSell.String() {
							t.tp.PushNewOrder(trading_core.NewAskMarketQtyItem(data.OrderId, d(data.Qty), data.At))
						} else if side == trading_core.OrderSideBuy.String() {
							t.tp.PushNewOrder(trading_core.NewBidMarketQtyItem(data.OrderId, d(data.Qty), d(data.MaxAmount), data.At))
						} else {
							logrus.Errorf("新订单参数错误: %s side只能是sell/buy", raw)
						}
					} else if order_type == "market_amount" {
						//按成交金额
						if side == trading_core.OrderSideSell.String() {
							t.tp.PushNewOrder(trading_core.NewAskMarketAmountItem(data.OrderId, d(data.Qty), d(data.MaxQty), data.At))
						} else if side == trading_core.OrderSideBuy.String() {
							t.tp.PushNewOrder(trading_core.NewBidMarketAmountItem(data.OrderId, d(data.Amount), data.At))
						} else {
							logrus.Errorf("新订单参数错误: %s side只能是sell/buy", raw)
						}
					}
				}

			}(raw)
		}()

	}
}
func (t *tengine) pull_cancel_order() {
	<-t.restore_done_signal

	key := types.FormatCancelOrder.Format(t.symbol)
	logrus.Infof("正在监听redis队列: %s", key)
	for {
		func() {
			rdc := app.RedisPool().Get()
			defer rdc.Close()

			if n, _ := redis.Int64(rdc.Do("Llen", key)); n == 0 || t.tp.IsPausePushNew() {
				time.Sleep(time.Duration(50) * time.Millisecond)
				return
			}

			raw, _ := redis.Bytes(rdc.Do("LPOP", key)) // rdc.LPop(cx, key).Bytes()

			var data cancel_order
			err := json.Unmarshal(raw, &data)
			if err != nil {
				logrus.Warnf("%s 解析json: %s 错误: %s", key, raw, err)
			}

			if data.OrderId != "" {
				logrus.Debugf("%s队列LPop: %s", key, raw)
				side := strings.ToLower(data.Side)
				if side == "ask" {
					t.tp.CancelOrder(trading_core.OrderSideSell, data.OrderId)
				} else if side == "bid" {
					t.tp.CancelOrder(trading_core.OrderSideBuy, data.OrderId)
				} else {
					logrus.Errorf("取消订单参数错误: %s 类型只能是ask/bid", raw)
				}
			}

		}()

	}
}
func (t *tengine) monitor_result() {
	<-t.restore_done_signal

	for {
		select {
		case data := <-t.tp.ChTradeResult:
			go func() {
				raw, _ := json.Marshal(data)
				t.push_match_result(raw)
			}()
		case uniq := <-t.tp.ChCancelResult:
			go func() {
				key := types.FormatCancelResult.Format(t.symbol)

				data := map[string]any{
					"order_id": uniq,
					"cancel":   "success",
				}

				rdc := app.RedisPool().Get()
				defer rdc.Close()

				raw, _ := json.Marshal(data)
				if _, err := rdc.Do("RPUSH", key, raw); err != nil { //rdc.RPush(cx, key, raw).Err()
					logrus.Warnf("%s队列RPush: %s %s", key, raw, err)
				}
			}()

		default:
			time.Sleep(time.Duration(50) * time.Millisecond)
		}

	}
}

func (t *tengine) push_match_result(data []byte) {
	rdc := app.RedisPool().Get()
	defer rdc.Close()

	key := types.FormatTradeResult.Format(t.symbol)
	if _, err := rdc.Do("RPUSH", key, data); err != nil {
		logrus.Warnf("往%s队列RPush: %s %s", key, data, err)
	}
	// if viper.GetBool("haotrader.notify_quote") {
	// 	quote_key := types.FormatQuoteTradeResult.Format(t.symbol)
	// 	rdc.RPush(cx, quote_key, data)
	// }
}
