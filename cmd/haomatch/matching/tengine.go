package matching

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/types"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

type tengine struct {
	symbol              string
	tp                  *trading_core.TradePair
	restore_done_signal chan struct{}
	update              chan struct{}
}

func NewTengine(symbol string, price_digit, qty_digit int) *trading_core.TradePair {
	te := tengine{
		symbol:              symbol,
		tp:                  trading_core.NewTradePair(symbol, price_digit, qty_digit),
		restore_done_signal: make(chan struct{}),
		update:              make(chan struct{}, 100),
	}

	go te.queue_monitor()
	go te.restore()
	go te.push_depth_to_redis()
	go te.pull_new_order()
	go te.pull_cancel_order()
	go te.monitor_result()
	return te.tp
}

func (t *tengine) push_depth_to_redis() {
	depth_topic := types.FormatDepthData.Format(t.symbol)

	//如果长时间没有触发，5s自动触发一次更新
	go func() {
		for {
			time.Sleep(time.Duration(5) * time.Second)
			t.update <- struct{}{}
		}
	}()

	for {
		select {
		case <-t.update:
			go func() {
				price, at := t.tp.LatestPrice()
				data := types.RedisDepthData{
					Price: price.String(),
					At:    at,
					Asks:  t.tp.GetAskDepth(50),
					Bids:  t.tp.GetBidDepth(50),
				}

				rdc := app.RedisPool().Get()
				defer rdc.Close()

				raw := data.JSON()
				if _, err := rdc.Do("SET", depth_topic, raw); err != nil {
					app.Logger.Errorf("set redis %s err: %s", depth_topic, err)
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
				Side:      qi.GetOrderSide(),
				OrderType: trading_core.OrderTypeLimit,
				Price:     qi.GetPrice().String(),
				Qty:       qi.GetQuantity().String(),
				At:        qi.GetCreateTime(),
			}

			if qi.GetQuantity().Cmp(decimal.Zero) > 0 {
				app.Logger.Debugf("queue event update: %#v", raw)
				go localdb.Set(t.symbol, raw.OrderId, raw.Json())
			}
		}
		t.update <- struct{}{}
	}, func(qi trading_core.QueueItem) {
		if t.tp.TriggerEvent() {
			raw := Order{
				OrderId:   qi.GetUniqueId(),
				Side:      qi.GetOrderSide(),
				OrderType: trading_core.OrderTypeLimit,
				Price:     qi.GetPrice().String(),
				Qty:       "0",
				At:        qi.GetCreateTime(),
			}
			app.Logger.Debugf("queue event remove: %#v", raw)
			go localdb.Remove(t.symbol, raw.OrderId)
		}

		t.update <- struct{}{}
	}, func(tl trading_core.TradeResult) {
		//只保留最近的1条成交记录,用于恢复最新成交价格
		localdb.Set("tradelog", t.symbol, tl.Json())
	})
}

func (t *tengine) restore() {

	defer func() {
		app.Logger.Infof("数据恢复 [%s] 已完成", t.symbol)
		close(t.restore_done_signal)
		t.tp.SetTriggerEvent(true)
		t.tp.SetPauseMatch(false)
	}()
	//从磁盘恢复上一次的数据，先暂停撮合系统的撮合，等数据全部加载完成后再开启撮合
	t.tp.SetPauseMatch(true)

	//恢复orderbook
	data := localdb.Find(t.symbol, "")
	app.Logger.Infof("正在恢复[%s]数据，共%d条", t.symbol, len(data))
	for i, v := range data {
		func(n int, raw []byte) {
			app.Logger.Infof("恢复数据[%s]第%d条: %s", t.symbol, n+1, raw)
			var data Order
			json.Unmarshal(raw, &data)

			if data.Side == trading_core.OrderSideSell {
				t.tp.PushNewOrder(trading_core.NewAskLimitItem(data.OrderId, utils.D(data.Price), utils.D(data.Qty), data.At))
			} else {
				t.tp.PushNewOrder(trading_core.NewBidLimitItem(data.OrderId, utils.D(data.Price), utils.D(data.Qty), data.At))
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
	app.Logger.Infof("监听队列: %s", key)

	rdc := app.RedisPool().Get()
	defer rdc.Close()

	for {
		func() {
			if n, _ := redis.Int64(rdc.Do("LLen", key)); n == 0 || t.tp.IsPausePushNew() {
				time.Sleep(time.Duration(50) * time.Millisecond)
				return
			}

			raw, err := redis.Bytes(rdc.Do("Lpop", key))
			if err != nil {
				app.Logger.Errorf("lpop %s err: %s", key, err.Error())
				return
			}

			go func(raw []byte) {
				var data Order
				err := json.Unmarshal(raw, &data)
				if err != nil {
					app.Logger.Warnf("%s 解析json: %s 错误: %s", key, raw, err)
				}

				if data.OrderId != "" {
					app.Logger.Debugf("收到新订单: %s", raw)
					// side := strings.ToLower(data.Side)
					// order_type := strings.ToLower(data.OrderType)

					if data.OrderType == trading_core.OrderTypeLimit {
						if data.Side == trading_core.OrderSideSell {
							t.tp.PushNewOrder(trading_core.NewAskLimitItem(data.OrderId, utils.D(data.Price), utils.D(data.Qty), data.At))
						} else if data.Side == trading_core.OrderSideBuy {
							t.tp.PushNewOrder(trading_core.NewBidLimitItem(data.OrderId, utils.D(data.Price), utils.D(data.Qty), data.At))
						} else {
							app.Logger.Errorf("新订单参数错误: %s side只能是sell/buy", raw)
						}
					} else if data.OrderType == trading_core.OrderTypeMarket {
						if utils.D(data.Qty).Cmp(utils.D("0")) > 0 {
							// 按成交量
							if data.Side == trading_core.OrderSideSell {
								t.tp.PushNewOrder(trading_core.NewAskMarketQtyItem(data.OrderId, utils.D(data.Qty), data.At))
							} else if data.Side == trading_core.OrderSideBuy {
								t.tp.PushNewOrder(trading_core.NewBidMarketQtyItem(data.OrderId, utils.D(data.Qty), utils.D(data.MaxAmount), data.At))
							} else {
								app.Logger.Errorf("新订单参数错误: %s side只能是sell/buy", raw)
							}
						} else if utils.D(data.Amount).Cmp(utils.D("0")) > 0 {
							//按成交金额
							if data.Side == trading_core.OrderSideSell {
								t.tp.PushNewOrder(trading_core.NewAskMarketAmountItem(data.OrderId, utils.D(data.Qty), utils.D(data.MaxQty), data.At))
							} else if data.Side == trading_core.OrderSideBuy {
								t.tp.PushNewOrder(trading_core.NewBidMarketAmountItem(data.OrderId, utils.D(data.Amount), data.At))
							} else {
								app.Logger.Errorf("新订单参数错误: %s side只能是sell/buy", raw)
							}
						} else {
							app.Logger.Warnf("市价订单参数错误: %s", raw)
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
	app.Logger.Infof("监听取消订单队列: %s", key)
	for {
		func() {
			rdc := app.RedisPool().Get()
			defer rdc.Close()

			if n, _ := redis.Int64(rdc.Do("Llen", key)); n == 0 || t.tp.IsPausePushNew() {
				time.Sleep(time.Duration(50) * time.Millisecond)
				return
			}

			raw, _ := redis.Bytes(rdc.Do("LPOP", key)) // rdc.LPop(cx, key).Bytes()

			var data StructCancelOrder
			err := json.Unmarshal(raw, &data)
			if err != nil {
				app.Logger.Warnf("%s 解析json: %s 错误: %s", key, raw, err)
			}

			if data.OrderId != "" {
				app.Logger.Debugf("收到取消订单: %s %s", key, raw)

				if data.Side == trading_core.OrderSideSell {
					t.tp.CancelOrder(trading_core.OrderSideSell, data.OrderId)
				} else if data.Side == trading_core.OrderSideBuy {
					t.tp.CancelOrder(trading_core.OrderSideBuy, data.OrderId)
				} else {
					app.Logger.Errorf("取消订单参数错误: %s 类型只能是ask/bid", raw)
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
				raw, err := json.Marshal(data)
				if err != nil {
					app.Logger.Warnf("log: %v %s", data, err.Error())
					return
				}
				t.push_match_result(raw)
			}()
		case uniq := <-t.tp.ChCancelResult:
			go func() {
				key := types.FormatCancelResult.Format(t.symbol)

				data := StructCancelOrderResult{
					OrderId: uniq,
					Status:  "success",
				}

				rdc := app.RedisPool().Get()
				defer rdc.Close()

				raw := data.Json()
				if _, err := rdc.Do("RPUSH", key, raw); err != nil { //rdc.RPush(cx, key, raw).Err()
					app.Logger.Warnf("%s队列RPush: %s %s", key, raw, err)
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
		app.Logger.Warnf("撮合结果推送 %s 失败: %s %s", key, data, err)
	}

}
