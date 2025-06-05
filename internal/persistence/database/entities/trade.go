package entities

import (
	"fmt"

	"github.com/shopspring/decimal"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
)

type TradeRecord struct {
	ID
	Symbol     string                 `gorm:"-" json:"-"` // 忽略字段
	TradeId    string                 `gorm:"type:varchar(30);uniqueIndex:trade;not null" json:"trade_id"`
	Ask        string                 `gorm:"type:varchar(30);uniqueIndex:trade;not null" json:"ask"`
	Bid        string                 `gorm:"type:varchar(30);uniqueIndex:trade;not null" json:"bid"`
	TradeBy    matching_types.TradeBy `gorm:"type:smallint;default:0" json:"trade_by"`
	AskUid     string                 `gorm:"not null" json:"ask_uid"`
	BidUid     string                 `gorm:"not null" json:"bid_uid"`
	Price      decimal.Decimal        `gorm:"type:decimal(40,20);not null;default:0" json:"price"`
	Quantity   decimal.Decimal        `gorm:"type:decimal(40,20);not null;default:0" json:"quantity"`
	Amount     decimal.Decimal        `gorm:"type:decimal(40,20);not null;default:0" json:"amount"`
	AskFeeRate decimal.Decimal        `gorm:"type:decimal(40,20);not null;default:0" json:"ask_fee_rate"`
	AskFee     decimal.Decimal        `gorm:"type:decimal(40,20);not null;default:0" json:"ask_fee"`
	BidFeeRate decimal.Decimal        `gorm:"type:decimal(40,20);not null;default:0" json:"bid_fee_rate"`
	BidFee     decimal.Decimal        `gorm:"type:decimal(40,20);not null;default:0" json:"bid_fee"`
	BaseAt
}

func (t *TradeRecord) TableName() string {
	return fmt.Sprintf("trade_record_%s", t.Symbol)
}
