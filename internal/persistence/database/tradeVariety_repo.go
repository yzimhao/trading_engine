package database

import (
	"context"
	"strings"

	k_repo "github.com/duolacloud/crud-core-gorm/repositories"
	"github.com/duolacloud/crud-core/cache"
	"github.com/duolacloud/crud-core/datasource"
	b_mappers "github.com/duolacloud/crud-core/mappers"
	"github.com/duolacloud/crud-core/repositories"
	models_variety "github.com/yzimhao/trading_engine/v2/internal/models/variety"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	"gorm.io/gorm"
)

type gormTradeVarietyRepo struct {
	*repositories.MapperRepository[models_variety.TradeVariety, models_variety.CreateTradeVariety, models_variety.UpdateTradeVariety, entities.TradeVariety, entities.TradeVariety, map[string]any]
	varietyRepo persistence.VarietyRepository
}

func NewTradeVarietyRepo(datasource datasource.DataSource[gorm.DB], cache cache.Cache, varietyRepo persistence.VarietyRepository) persistence.TradeVarietyRepository {
	cacheWrapperRepo := repositories.NewCacheRepository(
		k_repo.NewGormCrudRepository[entities.TradeVariety, entities.TradeVariety, map[string]any](datasource),
		cache,
	)

	mapperRepo := repositories.NewMapperRepository(
		cacheWrapperRepo,
		b_mappers.NewJSONMapper[models_variety.TradeVariety, models_variety.CreateTradeVariety, models_variety.UpdateTradeVariety, entities.TradeVariety, entities.TradeVariety, map[string]any](),
	)

	return &gormTradeVarietyRepo{
		MapperRepository: mapperRepo,
		varietyRepo:      varietyRepo,
	}
}

func (v *gormTradeVarietyRepo) FindBySymbol(ctx context.Context, symbol string) (tradeVariety *models_variety.TradeVariety, err error) {
	symbol = strings.ToLower(symbol)
	tradeVariety, err = v.QueryOne(ctx, map[string]any{
		"symbol": map[string]any{
			"eq": symbol,
		},
	})
	if err != nil {
		return nil, err
	}

	tradeVariety.BaseVariety, err = v.varietyRepo.QueryOne(ctx, map[string]any{
		"id": map[string]any{
			"eq": tradeVariety.BaseId,
		},
	})
	if err != nil {
		return nil, err
	}

	tradeVariety.TargetVariety, err = v.varietyRepo.QueryOne(ctx, map[string]any{
		"id": map[string]any{
			"eq": tradeVariety.TargetId,
		},
	})
	if err != nil {
		return nil, err
	}

	return tradeVariety, nil
}
