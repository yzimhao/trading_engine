package kline

import (
	"time"

	"github.com/yzimhao/trading_engine/v2/internal/models"
	kline_types "github.com/yzimhao/trading_engine/v2/pkg/kline/types"
)

type Kline struct {
	models.UUID
	models.Base
	Symbol  string                 `json:"-"`
	Period  kline_types.PeriodType `json:"-"`
	OpenAt  time.Time              `json:"open_at,omitempty"`
	CloseAt time.Time              `json:"close_at,omitempty"`
	Open    string                 `json:"open,omitempty"`
	High    string                 `json:"high,omitempty"`
	Low     string                 `json:"low,omitempty"`
	Close   string                 `json:"close,omitempty"`
	Volume  string                 `json:"volume,omitempty"`
	Amount  string                 `json:"amount,omitempty"`
}

type CreateKline struct {
	Symbol  string
	Period  kline_types.PeriodType
	OpenAt  time.Time `json:"open_at,omitempty"`
	CloseAt time.Time `json:"close_at,omitempty"`
	Open    string    `json:"open,omitempty"`
	High    string    `json:"high,omitempty"`
	Low     string    `json:"low,omitempty"`
	Close   string    `json:"close,omitempty"`
	Volume  string    `json:"volume,omitempty"`
	Amount  string    `json:"amount,omitempty"`
}

type UpdateKline struct {
	Symbol  string
	Period  kline_types.PeriodType
	OpenAt  time.Time `json:"open_at,omitempty"`
	CloseAt time.Time `json:"close_at,omitempty"`
	Open    *string   `json:"open,omitempty"`
	High    *string   `json:"high,omitempty"`
	Low     *string   `json:"low,omitempty"`
	Close   *string   `json:"close,omitempty"`
	Volume  *string   `json:"volume,omitempty"`
	Amount  *string   `json:"amount,omitempty"`
}
