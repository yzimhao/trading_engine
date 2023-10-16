package orders

import (
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

func (u *UnfinishedOrder) Create(db *xorm.Session) error {
	//todo 频繁查询表是否存在，后面考虑缓存一下
	exist, err := db.IsTableExist(u.TableName())
	if err != nil {
		return err
	}
	if !exist {
		err := db.CreateTable(u)
		if err != nil {
			return err
		}

		err = db.CreateIndexes(u)
		if err != nil {
			return err
		}

		err = db.CreateUniques(u)
		if err != nil {
			return err
		}
	}

	_, err = db.Insert(u)
	if err != nil {
		return err
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
