package assets

import (
	"fmt"

	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

// 用户资产变动记录
type AssetsLog struct {
	Id         int64      `xorm:"pk autoincr bigint" json:"id"`
	UserId     string     `xorm:"varchar(30) index notnull" json:"user_id"`
	Symbol     string     `xorm:"varchar(30) index notnull" json:"symbol"`
	Before     string     `xorm:"decimal(40,20) default(0)" json:"before"`                    // 变动前
	Amount     string     `xorm:"decimal(40,20) default(0)" json:"amount"`                    // 变动数
	After      string     `xorm:"decimal(40,20) default(0)" json:"after"`                     // 变动后
	BusinessId string     `xorm:"varchar(100) index(business_id) notnull" json:"business_id"` //业务相关的id
	OpType     OpBehavior `xorm:"varchar(15)" json:"op_type"`
	Info       string     `xorm:"varchar(64)" json:"info"`
	CreateTime utils.Time `xorm:"timestamp created" json:"create_time"`
	UpdateTime utils.Time `xorm:"timestamp updated" json:"update_time"`
}

func (a *AssetsLog) TableName() string {
	return fmt.Sprintf("assets_log_%s", a.Symbol)
}

func QueryAssetsLogBusIdIsExist(symbol string, user_id string, business_id string) bool {
	db := app.Database().NewSession()
	defer db.Close()

	table := AssetsLog{Symbol: symbol}
	ok, _ := db.Table(&table).Where("user_id=? and business_id=?", user_id, business_id).Exist()
	return ok
}
