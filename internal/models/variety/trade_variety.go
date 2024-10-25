package variety

import (
	"github.com/yzimhao/trading_engine/v2/internal/models"
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
)

type TradeVariety struct {
	ID             int32        `json:"id"`
	Symbol         string       `json:"symbol"`
	Name           string       `json:"name"`
	TargetId       int32        `json:"target_id"`
	BaseId         int32        `json:"base_id"`
	TargetVariety  *Variety     `json:"target"`
	BaseVariety    *Variety     `json:"base"`
	PriceDecimals  int          `json:"price_decimals"`
	QtyDecimals    int          `json:"qty_decimals"`
	AllowMinQty    string       `json:"allow_min_qty"`
	AllowMaxQty    string       `json:"allow_max_qty"`
	AllowMinAmount string       `json:"allow_min_amount"`
	AllowMaxAmount string       `json:"allow_max_amount"`
	FeeRate        string       `json:"fee_rate"`
	Status         types.Status `json:"status"`
	Sort           int64        `json:"sort"`
	models.Base
}

type CreateTradeVariety struct {
	Symbol         string       `json:"symbol"`
	Name           string       `json:"name"`
	TargetId       int32        `json:"target_id"`
	BaseId         int32        `json:"base_id"`
	PriceDecimals  int          `json:"price_decimals"`
	QtyDecimals    int          `json:"qty_decimals"`
	AllowMinQty    string       `json:"allow_min_qty"`
	AllowMaxQty    string       `json:"allow_max_qty"`
	AllowMinAmount string       `json:"allow_min_amount"`
	AllowMaxAmount string       `json:"allow_max_amount"`
	FeeRate        string       `json:"fee_rate"`
	Status         types.Status `json:"status"`
	Sort           int64        `json:"sort"`
}

type UpdateTradeVariety struct {
	ID             int32         `json:"id"`
	Symbol         *string       `json:"symbol"`
	Name           *string       `json:"name"`
	TargetId       *int32        `json:"target_id"`
	BaseId         *int32        `json:"base_id"`
	PriceDecimals  *int          `json:"price_decimals"`
	QtyDecimals    *int          `json:"qty_decimals"`
	AllowMinQty    *string       `json:"allow_min_qty"`
	AllowMaxQty    *string       `json:"allow_max_qty"`
	AllowMinAmount *string       `json:"allow_min_amount"`
	AllowMaxAmount *string       `json:"allow_max_amount"`
	FeeRate        *string       `json:"fee_rate"`
	Status         *types.Status `json:"status"`
	Sort           *int64        `json:"sort"`
}
