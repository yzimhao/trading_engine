package models

import "github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"

type AssetsFreeze struct {
	Base
	UserId       string                `json:"user_id"`
	Symbol       string                `json:"symbol"`
	Amount       string                `json:"amount"`
	FreezeAmount string                `json:"freeze_amount"`
	Status       entities.FreezeStatus `json:"status"`
	TransId      string                `json:"trans_id"`
	FreezeType   entities.FreezeType   `json:"freeze_type"`
	Info         string                `json:"info"`
}

type CreateAssetsFreeze struct {
	UserId       string                `json:"user_id"`
	Symbol       string                `json:"symbol"`
	Amount       string                `json:"amount"`
	FreezeAmount string                `json:"freeze_amount"`
	Status       entities.FreezeStatus `json:"status"`
	TransId      string                `json:"trans_id"`
	FreezeType   entities.FreezeType   `json:"freeze_type"`
	Info         string                `json:"info"`
}

type UpdateAssetsFreeze struct {
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
