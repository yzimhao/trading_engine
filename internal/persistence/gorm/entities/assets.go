package entities

type FreezeStatus int
type AssetsChangeType string

const (
	FreezeStatusNew          FreezeStatus     = 0
	FreezeStatusDone         FreezeStatus     = 1
	AssetsChangeTypeTrade    AssetsChangeType = "trade"
	AssetsChangeTypeRecharge AssetsChangeType = "deposit"
	AssetsChangeTypeWithdraw AssetsChangeType = "withdraw"
	AssetsChangeTypeTransfer AssetsChangeType = "transfer"
)

type Assets struct {
	UUID
	Base
	UserId    string `gorm:"type:varchar(30);not null;uniqueIndex:userid_symbol" json:"user_id"`
	Symbol    string `gorm:"type:varchar(30);not null;uniqueIndex:userid_symbol" json:"symbol"`
	Total     string `gorm:"type:decimal(40,20);default:0;not null" json:"total"`
	Freeze    string `gorm:"type:decimal(40,20);default:0;not null" json:"freeze"`
	Available string `gorm:"type:decimal(40,20);default:0;not null" json:"avail"`
}
