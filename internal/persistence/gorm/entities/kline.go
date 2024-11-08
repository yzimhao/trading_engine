package entities

import (
	"fmt"
	"time"

	kline_types "github.com/yzimhao/trading_engine/v2/pkg/kline/types"
)

type Kline struct {
	UUID
	Base
	Symbol  string
	Period  kline_types.PeriodType
	OpenAt  time.Time `gorm:"timestamp uniqueIndex: open_at" json:"open_at,omitempty"`
	CloseAt time.Time `gorm:"timestamp" json:"close_at,omitempty"`
	Open    string    `gorm:"type:decimal(40,20);default:0" json:"open,omitempty"`
	High    string    `gorm:"type:decimal(40,20);default:0" json:"high,omitempty"`
	Low     string    `gorm:"type:decimal(40,20);default:0" json:"low,omitempty"`
	Close   string    `gorm:"type:decimal(40,20);default:0" json:"close,omitempty"`
	Volume  string    `gorm:"type:decimal(40,20);default:0" json:"volume,omitempty"`
	Amount  string    `gorm:"type:decimal(40,20);default:0" json:"amount,omitempty"`
}

func (kl *Kline) TableName() string {
	return fmt.Sprintf("kline_%s_%s", kl.Symbol, kl.Period)
}
