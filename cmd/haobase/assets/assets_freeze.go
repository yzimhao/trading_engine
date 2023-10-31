package assets

import (
	"time"
)

type AssetsFreeze struct {
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
