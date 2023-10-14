package symbols

import (
	"time"

	"github.com/yzimhao/trading_engine/utils"
)

type status int

const (
	StatusDisable status = 0
	StatusEnabled status = 1
)

type Varieties struct {
	Id            int       `xorm:"pk autoincr int" json:"id"`
	Symbol        string    `xorm:"varchar(100) notnull unique(symbol)" json:"symbol"`
	Name          string    `xorm:"varchar(250) notnull" json:"name"`
	ShowPrecision int       `xorm:"default(0)" json:"show_precision"`
	MinPrecision  int       `xorm:"default(0)" json:"min_precision"`
	Base          bool      `xorm:"default(0)" json:"base"` //是否为本位币
	Status        status    `xorm:"default(0) notnull" json:"-"`
	CreateTime    time.Time `xorm:"timestamp created" json:"-"`
	UpdateTime    time.Time `xorm:"timestamp updated" json:"update_at"`
}

type TradingVarieties struct {
	Id             int               `xorm:"pk autoincr int" json:"-"`
	Symbol         string            `xorm:"varchar(100) notnull unique(symbol)" json:"symbol"`
	Name           string            `xorm:"varchar(250) notnull" json:"name"`
	TargetSymbolId int               `xorm:"default(0) unique(symbol_base)" json:"target_symbol_id"` //交易物品
	BaseSymbolId   int               `xorm:"default(0) unique(symbol_base)" json:"base_symbol_id"`   //支付货币
	PricePrecision int               `xorm:"default(2)" json:"price_precision"`
	QtyPrecision   int               `xorm:"default(0)" json:"qty_precision"`
	AllowMinQty    utils.FloatString `xorm:"decimal(40,20) default(0.01)" json:"allow_min_qty"`
	AllowMaxQty    utils.FloatString `xorm:"decimal(40,20) default(999999)" json:"allow_max_qty"`
	AllowMinAmount utils.FloatString `xorm:"decimal(40,20) default(0.01)" json:"allow_min_amount"`
	AllowMaxAmount utils.FloatString `xorm:"decimal(40,20) default(999999)" json:"allow_max_amount"`
	FeeRate        utils.FloatString `xorm:"decimal(40,20) default(0)" json:"fee_rate"`
	Status         status            `xorm:"default(0)" json:"-"`
	CreateTime     time.Time         `xorm:"timestamp created" json:"-"`
	UpdateTime     time.Time         `xorm:"timestamp updated" json:"update_at"`
	Target         Varieties         `xorm:"-" json:"target"`
	Base           Varieties         `xorm:"-" json:"Base"`
}
