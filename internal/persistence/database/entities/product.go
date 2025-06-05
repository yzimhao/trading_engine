package entities

import (
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
)

type Product struct {
	ID             int32        `gorm:"primaryKey;autoIncrement" json:"id"`
	Symbol         string       `gorm:"type:varchar(100); not null; uniqueIndex:symbol_idx" json:"symbol"`
	Name           string       `gorm:"type:varchar(250); not null" json:"name"`
	TargetId       int32        `gorm:"default:0; uniqueIndex:symbol_base_idx" json:"target_id"`
	BaseId         int32        `gorm:"default:0; uniqueIndex:symbol_base_idx" json:"base_id"`
	PriceDecimals  int          `gorm:"default:2" json:"price_decimals"`
	QtyDecimals    int          `gorm:"default:0" json:"qty_decimals"`
	AllowMinQty    string       `gorm:"type:decimal(40,20); default:0.01" json:"allow_min_qty"`
	AllowMaxQty    string       `gorm:"type:decimal(40,20); default:999999" json:"allow_max_qty"`
	AllowMinAmount string       `gorm:"type:decimal(40,20); default:0.01" json:"allow_min_amount"`
	AllowMaxAmount string       `gorm:"type:decimal(40,20); default:999999" json:"allow_max_amount"`
	FeeRate        string       `gorm:"type:decimal(40,20); default:0" json:"fee_rate"`
	Status         types.Status `gorm:"default:0" json:"status"`
	Sort           int64        `gorm:"default:0" json:"sort"`
	Base           *Asset       `gorm:"-" json:"base"`
	Target         *Asset       `gorm:"-" json:"target"`
	BaseAt
}
