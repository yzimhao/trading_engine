package entities

import (
	"fmt"

	models_types "github.com/yzimhao/trading_engine/v2/internal/models/types"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
)

type Order struct {
	UUID
	Base
	Symbol         string                   `gorm:"type:varchar(30)" json:"symbol"`
	OrderId        string                   `gorm:"type:varchar(30);uniqueIndex:order_id;not null" json:"order_id"`
	OrderSide      matching_types.OrderSide `gorm:"type:varchar(10);index:order_side" json:"order_side"`
	OrderType      matching_types.OrderType `gorm:"type:varchar(10)" json:"order_type"` //价格策略，市价单，限价单
	UserId         string                   `gorm:"index:userid;not null" json:"user_id"`
	Price          string                   `gorm:"type:decimal(40,20);not null;default:0" json:"price"`
	Quantity       string                   `gorm:"type:decimal(40,20);not null;default:0" json:"quantity"`
	FeeRate        string                   `gorm:"type:decimal(40,20);not null;default:0" json:"fee_rate"`
	Amount         string                   `gorm:"type:decimal(40,20);not null;default:0" json:"amount"`
	FreezeQty      string                   `gorm:"type:decimal(40,20);not null;default:0" json:"freeze_qty"`
	FreezeAmount   string                   `gorm:"type:decimal(40,20);not null;default:0" json:"freeze_amount"`
	AvgPrice       string                   `gorm:"type:decimal(40,20);not null;default:0" json:"avg_price"` //订单撮合成功 结算逻辑写入的字段
	FinishedQty    string                   `gorm:"type:decimal(40,20);not null;default:0" json:"finished_qty"`
	FinishedAmount string                   `gorm:"type:decimal(40,20);not null;default:0" json:"finished_amount"`
	Fee            string                   `gorm:"type:decimal(40,20);not null;default:0" json:"fee"`
	Status         models_types.OrderStatus `gorm:"type:tinyint(1);default:0" json:"status"`
	NanoTime       int64                    `gorm:"type:bigint(20);not null;default:0" json:"nano_time"`
}

func (o *Order) TableName() string {
	return fmt.Sprintf("order_%s", o.Symbol)
}

type UnfinishedOrder struct {
	Order
}

func (o *UnfinishedOrder) TableName() string {
	return "unfinished_order"
}
