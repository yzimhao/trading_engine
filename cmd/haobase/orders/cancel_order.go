package orders

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gookit/goutil/arrutil"
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haobase/message"
	"github.com/yzimhao/trading_engine/cmd/haobase/message/ws"
	"github.com/yzimhao/trading_engine/cmd/haomatch/matching"
	"github.com/yzimhao/trading_engine/types"
	"github.com/yzimhao/trading_engine/utils/app"
	"github.com/yzimhao/trading_engine/utils/app/config"
)

func SubmitOrderCancel(order_id string) error {
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
		}
		rdc.Do("rpush", types.FormatCancelOrder.Format(order.Symbol), cancel.Json())
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
			var data matching.StructCancelOrderResult
			json.Unmarshal(raw, &data)
			go cancel_order(symbol, data.OrderId, 0)
		}()

	}
}

func cancel_order(symbol, order_id string, retry int) {
	lock := GetLock(ClearingLock, order_id)
	wait := 10

	app.Logger.Infof("取消订单%s lock: %d retry: %d", order_id, lock, retry)
	if lock > 0 && retry <= wait*2 {
		//等待10s 还是有锁，记录下订单，退出取消逻辑
		time.Sleep(time.Duration(500) * time.Millisecond)
		cancel_order(symbol, order_id, retry+1)
		return
	}
	if lock > 0 && retry > wait*2 {
		app.Logger.Errorf("取消%s订单%s失败", symbol, order_id)
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
			app.Logger.Errorf("取消订单 %s 失败 %s", order_id, err.Error())
			db.Rollback()
		} else {
			if err := db.Commit(); err != nil {
				app.Logger.Errorf("取消订单 %s 失败 %s", order_id, err.Error())
			} else {
				//取消成功websocket发送消息给前端
				to := fmt.Sprintf("user.%s", item.UserId)
				message.Publish(ws.MsgBody{
					To: to,
					Response: ws.Response{
						Type: types.MsgUserOrderCancel.Format(symbol),
						Body: map[string]string{
							"order_id": order_id,
						},
					},
				})
			}
		}
	}()

	tablename := GetOrderTableName(symbol)
	_, err = db.Table(tablename).Where("order_id=?", order_id).Get(&item)
	if err != nil {
		return
	}
	_, err = db.Table(new(UnfinishedOrder)).Where("order_id=?", order_id).Delete()
	if err != nil {
		return
	}

	//更新订单状态
	item.Status = OrderStatusCancel
	_, err = db.Table(tablename).Where("order_id=?", order_id).Cols("status").Update(item)
	if err != nil {
		return
	}

	//解除订单冻结金额
	_, err = assets.UnfreezeAllAssets(db, item.UserId, item.OrderId)
	if err != nil {
		return
	}
}
