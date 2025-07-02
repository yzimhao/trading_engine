package entities

import (
	"github.com/yzimhao/trading_engine/v2/internal/types"
)

/**
 * 资产信息
 * 例如：BTC、USDT这类的产品基本信息
 */

type Asset struct {
	ID           int32        `gorm:"primaryKey;autoIncrement" json:"id"`
	Symbol       string       `gorm:"type:varchar(100);not null;uniqueIndex:symbol" json:"symbol"`
	Name         string       `gorm:"type:varchar(250);not null" json:"name"`
	ShowDecimals int          `gorm:"default(0)" json:"show_decimals"`
	MinDecimals  int          `gorm:"default(0)" json:"min_decimals"`
	IsBase       bool         `gorm:"default(0)" json:"is_base"`
	Sort         int64        `gorm:"default(0)" json:"sort"`
	Status       types.Status `gorm:"default(0) notnull" json:"status"`
	BaseAt
}
