package gorm

import (
	"context"

	k_repo "github.com/duolacloud/crud-core-gorm/repositories"
	"github.com/duolacloud/crud-core/cache"
	"github.com/duolacloud/crud-core/datasource"
	b_mappers "github.com/duolacloud/crud-core/mappers"
	"github.com/duolacloud/crud-core/repositories"
	models_variety "github.com/yzimhao/trading_engine/v2/internal/models/variety"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
	"gorm.io/gorm"
)

type gormVarietyRepo struct {
	*repositories.MapperRepository[models_variety.Variety, models_variety.CreateVariety, models_variety.UpdateVariety, entities.Variety, entities.Variety, map[string]any]
}

func NewVarietyRepo(datasource datasource.DataSource[gorm.DB], cache cache.Cache) persistence.VarietyRepository {
	cacheWrapperRepo := repositories.NewCacheRepository(
		k_repo.NewGormCrudRepository[entities.Variety, entities.Variety, map[string]any](datasource),
		cache,
	)

	mapperRepo := repositories.NewMapperRepository(
		cacheWrapperRepo,
		b_mappers.NewJSONMapper[models_variety.Variety, models_variety.CreateVariety, models_variety.UpdateVariety, entities.Variety, entities.Variety, map[string]any](),
	)

	return &gormVarietyRepo{
		MapperRepository: mapperRepo,
	}
}

type gormTradeVarietyRepo struct {
	*repositories.MapperRepository[models_variety.TradeVariety, models_variety.CreateTradeVariety, models_variety.UpdateTradeVariety, entities.TradeVariety, entities.TradeVariety, map[string]any]
}

func NewTradeVarietyRepo(datasource datasource.DataSource[gorm.DB], cache cache.Cache) persistence.TradeVarietyRepository {
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
	}
}

func (v *gormTradeVarietyRepo) FindBySymbol(ctx context.Context, symbol string) (tradeVariety *models_variety.TradeVariety, err error) {
	tradeVariety, err = v.QueryOne(ctx, map[string]any{
		"symbol": map[string]any{
			"eq": symbol,
		},
	})
	return
}
