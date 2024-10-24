package persistence

import (
	"context"

	"github.com/duolacloud/crud-core/repositories"
	models_variety "github.com/yzimhao/trading_engine/v2/internal/models/variety"
)

type VarietyRepository interface {
	repositories.CrudRepository[models_variety.Variety, models_variety.CreateVariety, models_variety.UpdateVariety]
}

type TradeVarietyRepository interface {
	repositories.CrudRepository[models_variety.TradeVariety, models_variety.CreateTradeVariety, models_variety.UpdateTradeVariety]
	FindBySymbol(ctx context.Context, symbol string) (tradeVariety *models_variety.TradeVariety, err error)
}
