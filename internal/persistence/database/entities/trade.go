package entities

import (
	"fmt"

	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
)

type TradeLog struct {
	UUID
	BaseAt
	Symbol     string                 `gorm:"-" json:"-"` // 忽略字段
	TradeId    string                 `gorm:"type:varchar(30);uniqueIndex:trade;not null" json:"trade_id"`
	Ask        string                 `gorm:"type:varchar(30);uniqueIndex:trade;not null" json:"ask"`
	Bid        string                 `gorm:"type:varchar(30);uniqueIndex:trade;not null" json:"bid"`
	TradeBy    matching_types.TradeBy `gorm:"type:smallint;default:0" json:"trade_by"`
	AskUid     string                 `gorm:"not null" json:"ask_uid"`
	BidUid     string                 `gorm:"not null" json:"bid_uid"`
	Price      string                 `gorm:"type:decimal(40,20);not null;default:0" json:"price"`
	Quantity   string                 `gorm:"type:decimal(40,20);not null;default:0" json:"quantity"`
	Amount     string                 `gorm:"type:decimal(40,20);not null;default:0" json:"amount"`
	AskFeeRate string                 `gorm:"type:decimal(40,20);not null;default:0" json:"ask_fee_rate"`
	AskFee     string                 `gorm:"type:decimal(40,20);not null;default:0" json:"ask_fee"`
	BidFeeRate string                 `gorm:"type:decimal(40,20);not null;default:0" json:"bid_fee_rate"`
	BidFee     string                 `gorm:"type:decimal(40,20);not null;default:0" json:"bid_fee"`
}

func (t *TradeLog) TableName() string {
	return fmt.Sprintf("trade_log_%s", t.Symbol)
}
