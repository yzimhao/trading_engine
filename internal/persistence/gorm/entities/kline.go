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
	Open    string    `gorm:"decimal(30, 10)" json:"open,omitempty"`
	High    string    `gorm:"decimal(30, 10)" json:"high,omitempty"`
	Low     string    `gorm:"decimal(30, 10)" json:"low,omitempty"`
	Close   string    `gorm:"decimal(30, 10)" json:"close,omitempty"`
	Volume  string    `gorm:"decimal(30, 10)" json:"volume,omitempty"`
	Amount  string    `gorm:"decimal(30, 10)" json:"amount,omitempty"`
}

func (kl *Kline) TableName() string {
	return fmt.Sprintf("kline_%s_%s", kl.Symbol, kl.Period)
}
