package varieties

import (
	"github.com/yzimhao/trading_engine/types"
	"github.com/yzimhao/trading_engine/utils"
)

type Varieties struct {
	Id            int          `xorm:"pk autoincr int" json:"id"`
	Symbol        string       `xorm:"varchar(100) notnull unique(symbol)" json:"symbol"`
	Name          string       `xorm:"varchar(250) notnull" json:"name"`
	ShowPrecision int          `xorm:"default(0)" json:"show_precision"`
	MinPrecision  int          `xorm:"default(0)" json:"min_precision"`
	Base          bool         `xorm:"default(0)" json:"base"` //是否为本位币
	Sort          int64        `xorm:"default(0)" json:"sort"`
	Status        types.Status `xorm:"default(0) notnull" json:"status"`
	CreateTime    utils.Time   `xorm:"timestamp created" json:"create_time"`
	UpdateTime    utils.Time   `xorm:"timestamp updated" json:"update_time"`
}

type TradingVarieties struct {
	Id             int           `xorm:"pk autoincr int" json:"id"`
	Symbol         string        `xorm:"varchar(100) notnull unique(symbol)" json:"symbol"`
	Name           string        `xorm:"varchar(250) notnull" json:"name"`
	TargetSymbolId int           `xorm:"default(0) unique(symbol_base)" json:"target_symbol_id"` //交易物品
	BaseSymbolId   int           `xorm:"default(0) unique(symbol_base)" json:"base_symbol_id"`   //支付货币
	PricePrecision int           `xorm:"default(2)" json:"price_precision"`
	QtyPrecision   int           `xorm:"default(0)" json:"qty_precision"`
	AllowMinQty    utils.DeciStr `xorm:"decimal(40,20) default(0.01)" json:"allow_min_qty"`
	AllowMaxQty    utils.DeciStr `xorm:"decimal(40,20) default(999999)" json:"allow_max_qty"`
	AllowMinAmount utils.DeciStr `xorm:"decimal(40,20) default(0.01)" json:"allow_min_amount"`
	AllowMaxAmount utils.DeciStr `xorm:"decimal(40,20) default(999999)" json:"allow_max_amount"`
	FeeRate        utils.DeciStr `xorm:"decimal(40,20) default(0)" json:"fee_rate"`
	Status         types.Status  `xorm:"default(0)" json:"status"`
	CreateTime     utils.Time    `xorm:"timestamp created" json:"create_time"`
	UpdateTime     utils.Time    `xorm:"timestamp updated" json:"update_time"`
	Target         Varieties     `xorm:"-" json:"target"`
	Base           Varieties     `xorm:"-" json:"base"`
}
