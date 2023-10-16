package assets

import (
	"time"

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
	Behavior_Recharge OpBehavior = "recharge"
	Behavior_Withdraw OpBehavior = "withdraw"
	Behavior_Transfer OpBehavior = "transfer"
)

// 用户资产余额表
type Assets struct {
	Id         int64      `xorm:"pk autoincr bigint" json:"-"`
	UserId     string     `xorm:"varchar(30) notnull unique(userid_symbol)" json:"-"`
	Symbol     string     `xorm:"varchar(30) notnull unique(userid_symbol)" json:"symbol"`
	Total      string     `xorm:"decimal(40,20) default(0) notnull" json:"total"`
	Freeze     string     `xorm:"decimal(40,20) default(0) notnull" json:"freeze"`
	Available  string     `xorm:"decimal(40,20) default(0) notnull" json:"avail"`
	CreateTime time.Time  `xorm:"timestamp created" json:"-"`
	UpdateTime utils.Time `xorm:"timestamp updated" json:"update_at"`
}

// 用户资产变动记录
type assetsLog struct {
	Id         int64      `xorm:"pk autoincr bigint"`
	UserId     string     `xorm:"varchar(30) index notnull"`
	Symbol     string     `xorm:"varchar(30) index notnull"`
	Before     string     `xorm:"decimal(40,20) default(0)"`               // 变动前
	Amount     string     `xorm:"decimal(40,20) default(0)"`               // 变动数
	After      string     `xorm:"decimal(40,20) default(0)"`               // 变动后
	BusinessId string     `xorm:"varchar(100) index(business_id) notnull"` //业务相关的id
	Behavior   OpBehavior `xorm:"varchar(15)"`
	Info       string     `xorm:"varchar(64)"`
	CreateTime time.Time  `xorm:"timestamp created"`
	UpdateTime time.Time  `xorm:"timestamp updated"`
}

type assetsFreeze struct {
	Id           int64        `xorm:"pk autoincr bigint"`
	UserId       string       `xorm:"varchar(30) index notnull"`
	Symbol       string       `xorm:"varchar(30) index notnull"`
	Amount       string       `xorm:"decimal(40,20) default(0) notnull"`        // 冻结总量
	FreezeAmount string       `xorm:"decimal(40,20) default(0) notnull"`        // 冻结着的量
	Status       FreezeStatus `xorm:"tinyint(1)"`                               // 状态 冻结中, 已解冻
	BusinessId   string       `xorm:"varchar(100) unique(business_id) notnull"` //业务相关的id
	Info         string       `xorm:"varchar(64)"`
	CreateTime   time.Time    `xorm:"timestamp created"`
	UpdateTime   time.Time    `xorm:"timestamp updated"`
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
