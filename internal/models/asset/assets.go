package asset

import (
	"github.com/yzimhao/trading_engine/v2/internal/models"
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
)

type Asset struct {
	models.UUID
	models.Base
	UserId        string       `json:"user_id,omitempty"`
	Symbol        string       `json:"symbol,omitempty"`
	TotalBalance  types.Amount `json:"total_balance,omitempty"`
	FreezeBalance types.Amount `json:"freeze_balance,omitempty"`
	AvailBalance  types.Amount `json:"avail_balance,omitempty"`
}

type CreateAsset struct {
	UserId        string        `json:"user_id,omitempty"`
	Symbol        string        `json:"symbol,omitempty"`
	TotalBalance  *types.Amount `json:"total_balance,omitempty"`
	FreezeBalance *types.Amount `json:"freeze_balance,omitempty"`
	AvailBalance  *types.Amount `json:"avail_balance,omitempty"`
}

type UpdateAsset struct {
	models.UUID
	UserId        *string       `json:"user_id,omitempty"`
	Symbol        *string       `json:"symbol,omitempty"`
	TotalBalance  *types.Amount `json:"total_balance,omitempty"`
	FreezeBalance *types.Amount `json:"freeze_balance,omitempty"`
	AvailBalance  *types.Amount `json:"avail_balance,omitempty"`
}
