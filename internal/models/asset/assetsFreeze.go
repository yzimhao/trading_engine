package asset

import (
	"github.com/yzimhao/trading_engine/v2/internal/models"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
)

type AssetFreeze struct {
	models.Base
	UserId       string                `json:"user_id"`
	Symbol       string                `json:"symbol"`
	Amount       string                `json:"amount"`
	FreezeAmount string                `json:"freeze_amount"`
	Status       entities.FreezeStatus `json:"status"`
	TransId      string                `json:"trans_id"`
	FreezeType   entities.FreezeType   `json:"freeze_type"`
	Info         string                `json:"info"`
}

type CreateAssetFreeze struct {
	UserId       string                `json:"user_id"`
	Symbol       string                `json:"symbol"`
	Amount       string                `json:"amount"`
	FreezeAmount string                `json:"freeze_amount"`
	Status       entities.FreezeStatus `json:"status"`
	TransId      string                `json:"trans_id"`
	FreezeType   entities.FreezeType   `json:"freeze_type"`
	Info         string                `json:"info"`
}

type UpdateAssetFreeze struct {
	ID           int64                  `json:"id"`
	UserId       *string                `json:"user_id"`
	Symbol       *string                `json:"symbol"`
	Amount       *string                `json:"amount"`
	FreezeAmount *string                `json:"freeze_amount"`
	Status       *entities.FreezeStatus `json:"status"`
	TransId      *string                `json:"trans_id"`
	FreezeType   *entities.FreezeType   `json:"freeze_type"`
	Info         *string                `json:"info"`
}
