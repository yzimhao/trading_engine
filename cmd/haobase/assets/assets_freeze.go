package assets

import (
	"fmt"

	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

type AssetsFreeze struct {
	Id           int64        `xorm:"pk autoincr bigint" json:"id"`
	UserId       string       `xorm:"varchar(30) index notnull" json:"user_id"`
	Symbol       string       `xorm:"varchar(30) index notnull" json:"symbol"`
	Amount       string       `xorm:"decimal(40,20) default(0) notnull" json:"amount"`             // 冻结总量
	FreezeAmount string       `xorm:"decimal(40,20) default(0) notnull" json:"freeze_amount"`      // 冻结着的量
	Status       FreezeStatus `xorm:"tinyint(1)" json:"status"`                                    // 状态 冻结中, 已解冻
	BusinessId   string       `xorm:"varchar(100) unique(business_id) notnull" json:"business_id"` //业务相关的id
	OpType       OpBehavior   `xorm:"varchar(15)" json:"op_type"`
	Info         string       `xorm:"varchar(200)" json:"info"`
	CreateTime   utils.Time   `xorm:"timestamp created" json:"create_time"`
	UpdateTime   utils.Time   `xorm:"timestamp updated" json:"update_time"`
}

func (a *AssetsFreeze) TableName() string {
	return fmt.Sprintf("%sassets_freeze_%s", app.TablePrefix(), a.Symbol)
}
