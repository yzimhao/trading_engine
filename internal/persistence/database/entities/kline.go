package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	kline_types "github.com/yzimhao/trading_engine/v2/pkg/kline/types"
)

type Kline struct {
	UUID
	BaseAt
	Symbol  string                 `gorm:"-"`
	Period  kline_types.PeriodType `gorm:"-"`
	OpenAt  time.Time              `gorm:"timestamp uniqueIndex: open_close_at" json:"open_at,omitempty"`
	CloseAt time.Time              `gorm:"timestamp uniqueIndex: open_close_at" json:"close_at,omitempty"`
	Open    decimal.Decimal        `gorm:"type:decimal(40,20);default:0" json:"open,omitempty"`
	High    decimal.Decimal        `gorm:"type:decimal(40,20);default:0" json:"high,omitempty"`
	Low     decimal.Decimal        `gorm:"type:decimal(40,20);default:0" json:"low,omitempty"`
	Close   decimal.Decimal        `gorm:"type:decimal(40,20);default:0" json:"close,omitempty"`
	Volume  decimal.Decimal        `gorm:"type:decimal(40,20);default:0" json:"volume,omitempty"`
	Amount  decimal.Decimal        `gorm:"type:decimal(40,20);default:0" json:"amount,omitempty"`
}

func (kl *Kline) TableName() string {
	return fmt.Sprintf("kline_%s_%s", strings.ToLower(kl.Symbol), kl.Period)
}
