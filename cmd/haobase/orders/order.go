package orders

import (
	"fmt"
	"time"

	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/trading_core"
	"xorm.io/xorm"
)

type orderStatus string

const (
	OrderStatusNew    orderStatus = "new"
	OrderStatusDone   orderStatus = "done"
	OrderStatusCancel orderStatus = "cancel"
)

// 委托记录表
type Order struct {
	Id             int64                  `xorm:"pk autoincr bigint" json:"-"`
	Symbol         string                 `xorm:"varchar(30)" json:"symbol"`
	OrderId        string                 `xorm:"varchar(30) unique(order_id) notnull" json:"order_id"`
	OrderSide      trading_core.OrderSide `xorm:"varchar(10) index(order_side)" json:"order_side"`
	OrderType      trading_core.OrderType `xorm:"varchar(10)" json:"order_type"` //价格策略，市价单，限价单
	UserId         string                 `xorm:"index(userid) notnull" json:"-"`
	Price          string                 `xorm:"decimal(40,20) notnull default(0)" json:"price"`
	Quantity       string                 `xorm:"decimal(40,20) notnull default(0)" json:"quantity"`
	FeeRate        string                 `xorm:"decimal(40,20) notnull default(0)" json:"-"`
	Amount         string                 `xorm:"decimal(40,20) notnull default(0)" json:"-"`
	FreezeQty      string                 `xorm:"decimal(40,20) notnull default(0)" json:"-"`
	FreezeAmount   string                 `xorm:"decimal(40,20) notnull default(0)" json:"-"`
	AvgPrice       string                 `xorm:"decimal(40,20) notnull default(0)" json:"-"` //订单撮合成功 结算逻辑写入的字段
	FinishedQty    string                 `xorm:"decimal(40,20) notnull default(0)" json:"finished_qty"`
	FinishedAmount string                 `xorm:"decimal(40,20) notnull default(0)" json:"finished_amount"`
	Fee            string                 `xorm:"decimal(40,20) notnull default(0)" json:"-"`
	Status         orderStatus            `xorm:"varchar(10)" json:"status"`
	CreateTime     int64                  `xorm:"bigint" json:"create_time"` //时间戳 精确到纳秒
	UpdateTime     time.Time              `xorm:"timestamp updated" json:"-"`
}

func (o *Order) Save(db *xorm.Session) error {
	//todo 频繁查询表是否存在，后面考虑缓存一下
	exist, err := db.IsTableExist(o.TableName())
	if err != nil {
		return err
	}
	if !exist {
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

	o.CreateTime = time.Now().UnixNano()
	_, err = db.Table(o).Insert(o)
	if err != nil {
		return err
	}
	return nil
}

func (o *Order) TableName() string {
	return GetOrderTableName(o.Symbol)
}

func GetOrderTableName(symbol string) string {
	return fmt.Sprintf("order_%s", symbol)
}

func Find(symbol string, order_id string) *Order {
	db := assets.DB().NewSession()
	defer db.Close()

	var row Order
	db.Table(GetOrderTableName(symbol)).Where("order_id=?", order_id).Get(&row)
	if row.Id > 0 {
		return &row
	}
	return nil
}
