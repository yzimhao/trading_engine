package orders

import (
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/utils/app"
)

func cancel_order(symbol, order_id string) (order *Order, err error) {
	db := app.Database().NewSession()
	defer db.Close()

	err = db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	//todo 订单锁

	app.Logger.Infof("[orders] %s cancel: %s", symbol, order_id)

	order = &Order{
		Symbol:  symbol,
		OrderId: order_id,
		Status:  OrderStatusCancel,
	}
	//更新订单状态
	_, err = db.Table(new(UnfinishedOrder)).Where("order_id=?", order_id).Delete()
	if err != nil {
		return nil, err
	}
	_, err = db.Table(GetOrderTableName(symbol)).Where("order_id=?", order_id).Cols("status").Update(order)
	if err != nil {
		return nil, err
	}

	_, err = db.Table(order.TableName()).Where("order_id=?", order_id).Get(order)
	if err != nil {
		return nil, err
	}

	//解除订单冻结金额
	_, err = assets.UnfreezeAllAssets(db, order.UserId, order.OrderId)
	if err != nil {
		return nil, err
	}
	return order, nil
}
