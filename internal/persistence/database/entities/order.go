package entities

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/v2/internal/types"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
)

type Order struct {
	UUID
	BaseAt
	Symbol         string                   `gorm:"type:varchar(30)" json:"symbol"`
	OrderId        string                   `gorm:"type:varchar(30);uniqueIndex;not null" json:"order_id"`
	OrderSide      matching_types.OrderSide `gorm:"type:varchar(10);index;not null" json:"order_side"`
	OrderType      matching_types.OrderType `gorm:"type:varchar(10);not null" json:"order_type"` //价格策略，市价单，限价单
	UserId         string                   `gorm:"type:varchar(64);index;not null" json:"user_id"`
	Price          decimal.Decimal          `gorm:"type:decimal(40,20);not null;default:0" json:"price"`
	Quantity       decimal.Decimal          `gorm:"type:decimal(40,20);not null;default:0" json:"quantity"`
	FeeRate        decimal.Decimal          `gorm:"type:decimal(40,20);not null;default:0" json:"fee_rate"`
	Amount         decimal.Decimal          `gorm:"type:decimal(40,20);not null;default:0" json:"amount"`
	FreezeQty      decimal.Decimal          `gorm:"type:decimal(40,20);not null;default:0" json:"freeze_qty"`
	FreezeAmount   decimal.Decimal          `gorm:"type:decimal(40,20);not null;default:0" json:"freeze_amount"`
	AvgPrice       decimal.Decimal          `gorm:"type:decimal(40,20);not null;default:0" json:"avg_price"` //订单撮合成功 结算逻辑写入的字段
	FinishedQty    decimal.Decimal          `gorm:"type:decimal(40,20);not null;default:0" json:"finished_qty"`
	FinishedAmount decimal.Decimal          `gorm:"type:decimal(40,20);not null;default:0" json:"finished_amount"`
	Fee            decimal.Decimal          `gorm:"type:decimal(40,20);not null;default:0" json:"fee"`
	Status         types.OrderStatus        `gorm:"type:smallint;default:0" json:"status"`
	NanoTime       int64                    `gorm:"type:bigint;not null;default:0" json:"nano_time"`
}

func (o *Order) TableName() string {
	return fmt.Sprintf("order_%s", strings.ToLower(o.Symbol))
}

type UnfinishedOrder struct {
	Order
}

func (o *UnfinishedOrder) TableName() string {
	return "order_unfinished"
}
