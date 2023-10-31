package assets

import (
	"time"

	"github.com/yzimhao/trading_engine/utils/app"
)

// 用户资产变动记录
type AssetsLog struct {
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

func QueryAssetsLogBusIdIsExist(user_id string, business_id string) bool {
	db := app.Database().NewSession()
	defer db.Close()

	ok, _ := db.Table(new(AssetsLog)).Where("user_id=? and business_id=?", user_id, business_id).Exist()
	return ok
}
