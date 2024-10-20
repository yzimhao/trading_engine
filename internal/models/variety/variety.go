package variety

import (
	"github.com/yzimhao/trading_engine/v2/internal/models"
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
)

type Variety struct {
	ID           int32        `json:"id"`
	Symbol       string       `json:"symbol"`
	Name         string       `json:"name"`
	ShowDecimals int          `json:"show_decimals"`
	MinDecimals  int          `json:"min_decimals"`
	IsBase       bool         `json:"is_base"` //是否为本位币
	Sort         int64        `json:"sort"`
	Status       types.Status `json:"status"`
	models.Base
}

type CreateVariety struct {
	Symbol       string       `json:"symbol"`
	Name         string       `json:"name"`
	ShowDecimals int          `json:"show_decimals"`
	MinDecimals  int          `json:"min_decimals"`
	IsBase       bool         `json:"is_base"` //是否为本位币
	Sort         int64        `json:"sort"`
	Status       types.Status `json:"status"`
}

type UpdateVariety struct {
	ID           int32         `json:"id"`
	Symbol       *string       `json:"symbol"`
	Name         *string       `json:"name"`
	ShowDecimals *int          `json:"show_decimals"`
	MinDecimals  *int          `json:"min_decimals"`
	IsBase       *bool         `json:"is_base"` //是否为本位币
	Sort         *int64        `json:"sort"`
	Status       *types.Status `json:"status"`
}
