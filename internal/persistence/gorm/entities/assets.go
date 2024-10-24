package entities

import (
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
)

type FreezeStatus int
type AssetChangeType string
type FreezeType string

const (
	FreezeStatusNew         FreezeStatus    = 0
	FreezeStatusDone        FreezeStatus    = 1
	AssetChangeTypeTrade    AssetChangeType = "trade"
	AssetChangeTypeRecharge AssetChangeType = "deposit"
	AssetChangeTypeWithdraw AssetChangeType = "withdraw"
	AssetChangeTypeTransfer AssetChangeType = "transfer"
	FreezeTypeWithdraw      FreezeType      = "withdraw"
	FreezeTypeTransfer      FreezeType      = "transfer"
	FreezeTypeTrade         FreezeType      = "trade"

	SYSTEM_USER_ID = "system"
)

type Asset struct {
	UUID
	Base
	UserId        string       `gorm:"type:varchar(30);not null;uniqueIndex:userid_symbol" json:"user_id"`
	Symbol        string       `gorm:"type:varchar(30);not null;uniqueIndex:userid_symbol" json:"symbol"`
	TotalBalance  types.Amount `gorm:"type:decimal(40,20);default:0;not null" json:"total_balance"`
	FreezeBalance types.Amount `gorm:"type:decimal(40,20);default:0;not null" json:"freeze_balance"`
	AvailBalance  types.Amount `gorm:"type:decimal(40,20);default:0;not null" json:"avail_balance"`
}
