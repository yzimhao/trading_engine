package orders

import (
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/types/dbtables"
	"github.com/yzimhao/trading_engine/utils/app"
	"xorm.io/xorm"
)

// 未完全成交的委托订单记录表
type UnfinishedOrder struct {
	Order Order `xorm:"extends"`
}

func (u *UnfinishedOrder) TableName() string {
	return "order_unfinished"
}

func (u *UnfinishedOrder) Save(db *xorm.Session) error {
	if _, err := db.Insert(u); err != nil {
		return err
	}
	return nil
}

func (o *UnfinishedOrder) AutoCreateTable() error {
	db := app.Database().NewSession()
	defer db.Close()

	if !dbtables.Exist(db, o.TableName()) {
		err := db.CreateTable(o)
		if err != nil {
			return err
		}

		err = db.CreateIndexes(o)
		if err != nil {
			return err
		}

		err = db.CreateUniques(o)
		if err != nil {
			return err
		}
	}
	return nil
}

func FindUnfinished(symbol string, order_id string) *Order {
	db := app.Database().NewSession()
	defer db.Close()

	var row Order
	db.Table(new(UnfinishedOrder)).Where("order_id=?", order_id).Get(&row)
	if row.Id > 0 {
		return &row
	}
	return nil
}

func find_user_unfinished_orders(user_id string, symbol string, side trading_core.OrderSide) []Order {
	db := app.Database().NewSession()
	defer db.Close()

	var rows []Order
	db.Table(new(UnfinishedOrder)).Where("user_id=? and symbol=? and order_side=?", user_id, symbol, side.String()).OrderBy("price asc").Find(&rows)
	return rows
}
