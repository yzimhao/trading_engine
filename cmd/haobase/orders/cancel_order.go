package orders

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gookit/goutil/arrutil"
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haobase/base/varieties"
	"github.com/yzimhao/trading_engine/cmd/haobase/message"
	"github.com/yzimhao/trading_engine/cmd/haobase/message/ws"
	"github.com/yzimhao/trading_engine/cmd/haomatch/matching"
	"github.com/yzimhao/trading_engine/config"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/types"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func SubmitOrderCancel(order_id string, reason trading_core.CancelType) error {
	rdc := app.RedisPool().Get()
	defer rdc.Close()

	db := app.Database().NewSession()
	defer db.Close()

	var order Order
	has, _ := db.Table(new(UnfinishedOrder)).Where("order_id=?", order_id).Get(&order)
	if has {
		cancel := matching.StructCancelOrder{
			Side:    order.OrderSide,
			OrderId: order.OrderId,
			Reason:  reason,
		}
		rdc.Do("rpush", types.FormatCancelOrder.Format(order.Symbol), cancel.Json())
	} else {
		//已经完成或者已经被取消
		return fmt.Errorf("已经被取消或已完成")
	}

	return nil
}

func Run() {
	//监听取消订单队列
	local_config_symbols := config.App.Local.Symbols
	db_symbols := base.NewTSymbols().All()
	for _, item := range db_symbols {
		if len(local_config_symbols) > 0 && arrutil.Contains(local_config_symbols, item.Symbol) || len(local_config_symbols) == 0 {
			go watch_cancel_order_list(item.Symbol)
		}
	}

}

func watch_cancel_order_list(symbol string) {
	key := types.FormatCancelResult.Format(symbol)
	app.Logger.Infof("监听%s取消订单队列...", symbol)
	for {
		func() {
			rdc := app.RedisPool().Get()
			defer rdc.Close()

			if n, _ := redis.Int64(rdc.Do("LLen", key)); n == 0 {
				time.Sleep(time.Duration(50) * time.Millisecond)
				return
			}

			raw, _ := redis.Bytes(rdc.Do("Lpop", key))
			app.Logger.Infof("收到 %s 取消订单: %s", symbol, raw)
			var data matching.StructCancelOrder
			json.Unmarshal(raw, &data)
			go cancel_order(symbol, data, 0)
		}()

	}
}

func cancel_order(symbol string, cancel matching.StructCancelOrder, retry int) {
	lock := GetLock(SettleLock, cancel.OrderId)
	wait := 10

	app.Logger.Infof("取消订单%s lock: %d retry: %d", cancel.OrderId, lock, retry)
	if lock > 0 && retry <= wait*2 {
		//等待10s 还是有锁，记录下订单，退出取消逻辑
		time.Sleep(time.Duration(500) * time.Millisecond)
		cancel_order(symbol, cancel, retry+1)
		return
	}
	if lock > 0 && retry > wait*2 {
		app.Logger.Errorf("取消%s订单%s失败", symbol, cancel.OrderId)
		return
	}

	db := app.Database().NewSession()
	defer db.Close()

	var item Order

	err := db.Begin()
	if err != nil {
		app.Logger.Errorf(err.Error())
		return
	}

	defer func() {
		if err != nil {
			app.Logger.Errorf("取消订单 %s 失败 %s", cancel.OrderId, err.Error())
			db.Rollback()
		} else {
			if err := db.Commit(); err != nil {
				app.Logger.Errorf("取消订单 %s 失败 %s", cancel.OrderId, err.Error())
			} else {
				//取消成功websocket发送消息给前端
				to := types.MsgUser.Format(map[string]string{"user_id": item.UserId})
				message.Publish(ws.MsgBody{
					To: to,
					Response: ws.Response{
						Type: types.MsgOrderCancel.Format(map[string]string{"symbol": symbol}),
						Body: map[string]string{
							"order_id": cancel.OrderId,
						},
					},
				})
			}
		}
	}()

	tablename := &Order{Symbol: symbol}
	_, err = db.Table(tablename).Where("order_id=?", cancel.OrderId).Get(&item)
	if err != nil {
		return
	}
	_, err = db.Table(new(UnfinishedOrder)).Where("order_id=?", cancel.OrderId).Delete()
	if err != nil {
		return
	}

	//更新订单状态
	item.Status = OrderStatusCanceled
	if utils.D(item.Fee).Cmp(utils.D("0")) > 0 {
		if cancel.Reason == trading_core.CancelTypeByUser {
			item.Status = OrderStatusPartialCancel
		} else if cancel.Reason == trading_core.CancelTypeBySystem {
			item.Status = OrderStatusFilled
		}
	}

	_, err = db.Table(tablename).Where("order_id=?", cancel.OrderId).Cols("status").Update(item)
	if err != nil {
		return
	}

	varieties := varieties.NewTradingVarieties(symbol)
	//解除订单冻结的资产
	if item.OrderSide == trading_core.OrderSideSell {
		_, err = assets.UnfreezeAllAssets(db, varieties.Target.Symbol, item.UserId, item.OrderId)
		if err != nil {
			return
		}
	} else {
		_, err = assets.UnfreezeAllAssets(db, varieties.Base.Symbol, item.UserId, item.OrderId)
		if err != nil {
			return
		}
	}
}
