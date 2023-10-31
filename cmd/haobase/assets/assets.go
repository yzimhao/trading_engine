package assets

import (
	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

// 用户资产冻结记录
type FreezeStatus int
type OpBehavior string

const (
	FreezeStatusNew  FreezeStatus = 0
	FreezeStatusDone FreezeStatus = 1

	Behavior_Trade    OpBehavior = "trade"
	Behavior_Recharge OpBehavior = "deposit"
	Behavior_Withdraw OpBehavior = "withdraw"
	Behavior_Transfer OpBehavior = "transfer"
)

// 用户资产余额表
type Assets struct {
	Id         int64      `xorm:"pk autoincr bigint" json:"id"`
	UserId     string     `xorm:"varchar(30) notnull unique(userid_symbol)" json:"user_id"`
	Symbol     string     `xorm:"varchar(30) notnull unique(userid_symbol)" json:"symbol"`
	Total      string     `xorm:"decimal(40,20) default(0) notnull" json:"total"`
	Freeze     string     `xorm:"decimal(40,20) default(0) notnull" json:"freeze"`
	Available  string     `xorm:"decimal(40,20) default(0) notnull" json:"avail"`
	CreateTime utils.Time `xorm:"timestamp created" json:"create_time"`
	UpdateTime utils.Time `xorm:"timestamp updated" json:"update_time"`
}

func FindSymbol(user_id string, symbol string) *Assets {
	db := app.Database().NewSession()
	defer db.Close()

	var row Assets
	db.Table(new(Assets)).Where("user_id=? and symbol=?", user_id, symbol).Get(&row)
	return &row
}

func BalanceOfTotal(user_id, symbol string) decimal.Decimal {
	row := FindSymbol(user_id, symbol)
	if row.Id > 0 {
		return utils.D(row.Total)
	}
	return decimal.Zero
}

func BalanceOfFreeze(user_id, symbol string) decimal.Decimal {
	row := FindSymbol(user_id, symbol)
	if row.Id > 0 {
		return utils.D(row.Freeze)
	}
	return decimal.Zero
}

func BalanceOfAvailable(user_id, symbol string) decimal.Decimal {
	row := FindSymbol(user_id, symbol)
	if row.Id > 0 {
		return utils.D(row.Available)
	}
	return decimal.Zero
}
