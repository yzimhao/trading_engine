package settle

import (
	"fmt"

	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haobase/base/varieties"
	"github.com/yzimhao/trading_engine/cmd/haobase/orders"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/types/dbtables"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
	"xorm.io/xorm"
)

type clean struct {
	db                *xorm.Session
	trading_varieties *varieties.TradingVarieties
	ask               orders.Order
	bid               orders.Order
	tlog              trading_core.TradeResult
	tradelog          orders.TradeLog
	err               error
}

func newClean(raw trading_core.TradeResult) error {
	db := app.Database().NewSession()
	defer db.Close()

	tv, err := base.NewTSymbols().Get(raw.Symbol)
	if err != nil {
		app.Logger.Errorf("tsymbol error: %s", err)
		return err
	}

	//提前创建需要的表, 创建表的操作不能放到事务中
	dbtables.AutoCreateTable(db, &orders.TradeLog{Symbol: raw.Symbol})
	dbtables.AutoCreateTable(db, &orders.Order{Symbol: raw.Symbol})
	dbtables.AutoCreateTable(db, &assets.AssetsFreeze{Symbol: tv.Target.Symbol})
	dbtables.AutoCreateTable(db, &assets.AssetsFreeze{Symbol: tv.Base.Symbol})
	dbtables.AutoCreateTable(db, &assets.AssetsLog{Symbol: tv.Target.Symbol})
	dbtables.AutoCreateTable(db, &assets.AssetsLog{Symbol: tv.Base.Symbol})

	//

	item := clean{
		db:                db,
		trading_varieties: tv,
		ask:               orders.Order{},
		bid:               orders.Order{},
		tlog:              raw,
	}

	err = item.flow()

	//解锁
	orders.UnLock(orders.SettleLock, item.ask.OrderId)
	orders.UnLock(orders.SettleLock, item.bid.OrderId)

	//记录失败的订单
	if err != nil {
		app.Logger.Errorf("结算失败: %s %s ask:%+v bid:%+v", raw.Json(), err.Error(), item.ask, item.bid)
	} else {
		notify_quote(raw)
	}
	return err
}

func (c *clean) flow() error {

	c.db.Begin()
	defer func() {
		if c.err != nil {
			if err := c.db.Rollback(); err != nil {
				app.Logger.Errorf("结算事务回滚失败: %s", err.Error())
			}
		} else {
			if err := c.db.Commit(); err != nil {
				app.Logger.Errorf("结算事务提交失败: %s", err.Error())
			}
		}
	}()

	c.err = c.check_order()
	c.err = c.trade_log()
	c.err = c.update_order(trading_core.OrderSideSell)
	c.err = c.update_order(trading_core.OrderSideBuy)

	c.err = c.transfer()
	return c.err
}

func (c *clean) check_order() error {
	table := orders.Order{Symbol: c.tlog.Symbol}
	_, err := c.db.Table(&table).Where("order_id=?", c.tlog.AskOrderId).ForUpdate().Get(&c.ask)
	if err != nil {
		return err
	}
	_, err = c.db.Table(&table).Where("order_id=?", c.tlog.BidOrderId).ForUpdate().Get(&c.bid)
	if err != nil {
		return err
	}

	if c.ask.Status != orders.OrderStatusNew {
		return fmt.Errorf("卖单状态错误")
	}

	if c.bid.Status != orders.OrderStatusNew {
		return fmt.Errorf("买单状态错误")
	}
	return nil
}

func (c *clean) trade_log() error {
	amount := c.tlog.TradePrice.Mul(c.tlog.TradeQuantity)
	c.tradelog = orders.TradeLog{
		Symbol:  c.tlog.Symbol,
		TradeId: generate_trading_id(c.ask.OrderId, c.bid.OrderId),
		Ask:     c.ask.OrderId,
		Bid:     c.bid.OrderId,
		// TradeBy:    "",
		AskUid:     c.ask.UserId,
		BidUid:     c.bid.UserId,
		Price:      c.tlog.TradePrice.String(),
		Quantity:   c.tlog.TradeQuantity.String(),
		Amount:     amount.String(),
		AskFeeRate: c.ask.FeeRate,
		AskFee:     amount.Mul(utils.D(c.ask.FeeRate)).String(),
		BidFeeRate: c.bid.FeeRate,
		BidFee:     amount.Mul(utils.D(c.bid.FeeRate)).String(),
	}

	return c.tradelog.Save(c.db)
}

func (c *clean) update_order(side trading_core.OrderSide) error {
	var order *orders.Order
	if side == trading_core.OrderSideSell {
		order = &c.ask
		order.Fee = utils.D(order.Fee).Add(utils.D(c.tradelog.AskFee)).String()
	} else {
		order = &c.bid
		order.Fee = utils.D(order.Fee).Add(utils.D(c.tradelog.BidFee)).String()
	}

	order.Symbol = c.tlog.Symbol
	order.FinishedQty = utils.D(order.FinishedQty).Add(c.tlog.TradeQuantity).String()
	order.FinishedAmount = utils.D(order.FinishedAmount).Add(utils.D(c.tradelog.Amount)).String()
	order.AvgPrice = utils.D(order.FinishedAmount).Div(utils.D(order.FinishedQty)).String()

	if order.OrderType == trading_core.OrderTypeLimit {
		be := utils.D(order.FinishedQty).Cmp(utils.D(order.Quantity))
		if be > 0 {
			return fmt.Errorf("订单结算错误")
		}
		if be == 0 {
			order.Status = orders.OrderStatusDone
		}

		_, err := c.db.Table(order.TableName()).Where("order_id=?", order.OrderId).AllCols().Update(order)
		if err != nil {
			return err
		}

		if order.Status == orders.OrderStatusNew {
			_, err := c.db.Table(new(orders.UnfinishedOrder)).Where("order_id=?", order.OrderId).AllCols().Update(order)
			if err != nil {
				return err
			}
		} else {
			_, err := c.db.Table(new(orders.UnfinishedOrder)).Where("order_id=?", order.OrderId).Delete()
			if err != nil {
				return err
			}
		}
	} else {
		//市价单结算
		if utils.D(order.Quantity).Equal(utils.D(order.FinishedQty)) || utils.D(order.Amount).Equal(utils.D(order.FinishedAmount)) {
			order.Status = orders.OrderStatusDone
		}

		if c.tlog.Last == order.OrderId {
			order.Status = orders.OrderStatusDone
		}

		_, err := c.db.Table(order.TableName()).Where("order_id=?", order.OrderId).AllCols().Update(order)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *clean) transfer() error {
	//给买家结算交易物品
	_, err := assets.UnfreezeAssets(c.db, c.trading_varieties.Target.Symbol, c.ask.UserId, c.ask.OrderId, c.tlog.TradeQuantity.String())
	if err != nil {
		app.Logger.Errorf("解冻失败: %s %s", c.ask.OrderId, err.Error())
		return err
	}
	_, err = assets.Transfer(c.db, c.ask.UserId, c.bid.UserId, c.trading_varieties.Target.Symbol, c.tlog.TradeQuantity.String(), c.tradelog.TradeId, assets.Behavior_Trade)
	if err != nil {
		app.Logger.Errorf("Transfer: %s %s", c.trading_varieties.Target.Symbol, err.Error())
		return err
	}

	//卖家结算本位币
	amount := utils.D(c.tradelog.Amount).Add(utils.D(c.tradelog.BidFee))
	_, err = assets.UnfreezeAssets(c.db, c.trading_varieties.Base.Symbol, c.bid.UserId, c.bid.OrderId, amount.String())
	if err != nil {
		app.Logger.Errorf("解冻失败: %s %s", c.bid.OrderId, err.Error())
		return err
	}

	//扣除fee
	fee := utils.D(c.tradelog.BidFee).Add(utils.D(c.tradelog.AskFee))
	_, err = assets.Transfer(c.db, c.bid.UserId, c.ask.UserId, c.trading_varieties.Base.Symbol, amount.Sub(fee).String(), c.tradelog.TradeId, assets.Behavior_Trade)
	if err != nil {
		app.Logger.Errorf("Transfer: %s %s", c.trading_varieties.Base.Symbol, err.Error())
		return err
	}

	//手续费收入到一个全局的账号里
	_, err = assets.Transfer(c.db, c.bid.UserId, assets.UserSystemFee, c.trading_varieties.Base.Symbol, fee.String(), c.tradelog.TradeId, assets.Behavior_Trade)
	if err != nil {
		return err
	}

	// //市价单解除全部冻结
	// if c.tlog.Last != "" {
	// 	app.Logger.Infof("市价订单%s完成 解除剩余全部资产", c.tlog.Last)
	// 	if c.ask.OrderType == trading_core.OrderTypeMarket {
	// 		_, err = assets.UnfreezeAllAssets(c.db, c.trading_varieties.Target.Symbol, c.ask.UserId, c.ask.OrderId)
	// 		if err != nil {
	// 			app.Logger.Errorf("解冻UnfreezeAllAssets: %s %s", c.ask.OrderId, err.Error())
	// 			return err
	// 		}
	// 	}
	// 	if c.bid.OrderType == trading_core.OrderTypeMarket {
	// 		_, err = assets.UnfreezeAllAssets(c.db, c.trading_varieties.Base.Symbol, c.bid.UserId, c.bid.OrderId)
	// 		if err != nil {
	// 			app.Logger.Errorf("解冻UnfreezeAllAssets: %s %s", c.bid.OrderId, err.Error())
	// 			return err
	// 		}
	// 	}
	// }

	if c.ask.Status == orders.OrderStatusDone {
		_, err = assets.UnfreezeAllAssets(c.db, c.trading_varieties.Target.Symbol, c.ask.UserId, c.ask.OrderId)
		if err != nil {
			app.Logger.Errorf("解冻UnfreezeAllAssets: %s %s", c.ask.OrderId, err.Error())
			return err
		}
	}
	if c.bid.Status == orders.OrderStatusDone {
		_, err = assets.UnfreezeAllAssets(c.db, c.trading_varieties.Base.Symbol, c.bid.UserId, c.bid.OrderId)
		if err != nil {
			app.Logger.Errorf("解冻UnfreezeAllAssets: %s %s", c.bid.OrderId, err.Error())
			return err
		}
	}

	return nil
}
