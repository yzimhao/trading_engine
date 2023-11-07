package orders

import (
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/cmd/haobase/base/varieties"
	"github.com/yzimhao/trading_engine/cmd/haomatch/matching"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func NewMarketOrderByQty(user_id string, symbol string, side trading_core.OrderSide, qty string) (*Order, error) {
	return market_order_qty(user_id, symbol, side, qty)
}

func market_order_qty(user_id string, symbol string, side trading_core.OrderSide, qty string) (order *Order, err error) {
	varieties := varieties.NewTradingVarieties(symbol)

	neworder := Order{
		OrderId:        generate_order_id_by_side(side),
		Symbol:         symbol,
		OrderSide:      side,
		OrderType:      trading_core.OrderTypeMarket,
		UserId:         user_id,
		Price:          "0",
		AvgPrice:       "0",
		Quantity:       qty,
		FinishedQty:    "0",
		Fee:            "0",
		Amount:         "0",
		FreezeQty:      "0",
		FreezeAmount:   "0",
		FinishedAmount: "0",
		FeeRate:        string(varieties.FeeRate),
		Status:         OrderStatusNew,
	}

	if _, err := order_pre_inspection(varieties, &neworder); err != nil {
		return nil, err
	}

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

	//冻结资产
	if neworder.OrderSide == trading_core.OrderSideSell {
		_, err = assets.FreezeAssets(db, user_id, varieties.Target.Symbol, neworder.Quantity, neworder.OrderId, assets.Behavior_Trade)
		if err != nil {
			return nil, err
		}
		neworder.FreezeQty = neworder.Quantity
	} else if neworder.OrderSide == trading_core.OrderSideBuy {
		//冻结所有可用
		_, err = assets.FreezeTotalAssets(db, user_id, varieties.Base.Symbol, neworder.OrderId, assets.Behavior_Trade)
		if err != nil {
			return nil, err
		}

		freeze, err := assets.QueryFreeze(db, neworder.OrderId)
		if err != nil {
			return nil, err
		}
		neworder.FreezeAmount = freeze.FreezeAmount
	}

	if err = neworder.Save(db); err != nil {
		return nil, err
	}

	push_new_order_to_redis(neworder.Symbol, func() []byte {
		data := matching.Order{
			OrderId:   neworder.OrderId,
			OrderType: neworder.OrderType,
			Side:      neworder.OrderSide,
			Qty:       neworder.Quantity,
			MaxQty:    neworder.FreezeQty,
			Amount:    neworder.Amount,
			MaxAmount: neworder.FreezeAmount,
			At:        neworder.CreateTime,
		}
		return data.Json()
	}())

	return &neworder, nil
}

func NewMarketOrderByAmount(user_id string, symbol string, side trading_core.OrderSide, amount string) (order *Order, err error) {
	return market_order_amount(user_id, symbol, side, amount)
}

func market_order_amount(user_id string, symbol string, side trading_core.OrderSide, amount string) (order *Order, err error) {
	varieties := varieties.NewTradingVarieties(symbol)

	neworder := Order{
		OrderId:        generate_order_id_by_side(side),
		Symbol:         symbol,
		OrderSide:      side,
		OrderType:      trading_core.OrderTypeMarket,
		UserId:         user_id,
		Price:          "0",
		AvgPrice:       "0",
		Quantity:       "0",
		FinishedQty:    "0",
		Fee:            "0",
		FinishedAmount: "0",
		Amount:         amount,
		FreezeQty:      "0",
		FreezeAmount:   "0",
		FeeRate:        string(varieties.FeeRate),
		Status:         OrderStatusNew,
	}

	if _, err := order_pre_inspection(varieties, &neworder); err != nil {
		return nil, err
	}

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

	if neworder.OrderSide == trading_core.OrderSideSell {
		_, err = assets.FreezeTotalAssets(db, user_id, varieties.Target.Symbol, neworder.OrderId, assets.Behavior_Trade)
		if err != nil {
			return nil, err
		}

		freeze, err := assets.QueryFreeze(db, neworder.OrderId)
		if err != nil {
			return nil, err
		}
		neworder.FreezeQty = freeze.FreezeAmount

	} else if neworder.OrderSide == trading_core.OrderSideBuy {
		_, err = assets.FreezeAssets(db, user_id, varieties.Base.Symbol, neworder.Amount, neworder.OrderId, assets.Behavior_Trade)
		if err != nil {
			return nil, err
		}
		neworder.FreezeAmount = neworder.Amount
	}

	if err = neworder.Save(db); err != nil {
		return nil, err
	}

	push_new_order_to_redis(neworder.Symbol, func() []byte {
		fee := utils.D(neworder.FreezeAmount).Mul(utils.D(neworder.FeeRate))
		maxAmount := utils.D(neworder.FreezeAmount).Sub(fee)
		data := matching.Order{
			OrderId:   neworder.OrderId,
			OrderType: neworder.OrderType,
			Side:      neworder.OrderSide,
			Amount:    maxAmount.String(),
			MaxQty:    neworder.FreezeQty,
			MaxAmount: maxAmount.String(),
			At:        neworder.CreateTime,
		}
		return data.Json()
	}())

	return &neworder, nil
}
