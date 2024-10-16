package models

import (
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
)

type Assets struct {
	UUID
	Base
	UserId        string       `json:"user_id,omitempty"`
	Symbol        string       `json:"symbol,omitempty"`
	TotalBalance  types.Amount `json:"total_balance,omitempty"`
	FreezeBalance types.Amount `json:"freeze_balance,omitempty"`
	AvailBalance  types.Amount `json:"avail_balance,omitempty"`
}

type CreateAssets struct {
	UserId        string        `json:"user_id,omitempty"`
	Symbol        string        `json:"symbol,omitempty"`
	TotalBalance  *types.Amount `json:"total_balance,omitempty"`
	FreezeBalance *types.Amount `json:"freeze_balance,omitempty"`
	AvailBalance  *types.Amount `json:"avail_balance,omitempty"`
}

type UpdateAssets struct {
	UUID
	UserId        *string       `json:"user_id,omitempty"`
	Symbol        *string       `json:"symbol,omitempty"`
	TotalBalance  *types.Amount `json:"total_balance,omitempty"`
	FreezeBalance *types.Amount `json:"freeze_balance,omitempty"`
	AvailBalance  *types.Amount `json:"avail_balance,omitempty"`
}
