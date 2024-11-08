package gorm

import (
	k_repo "github.com/duolacloud/crud-core-gorm/repositories"
	"github.com/duolacloud/crud-core/cache"
	"github.com/duolacloud/crud-core/datasource"
	b_mappers "github.com/duolacloud/crud-core/mappers"
	"github.com/duolacloud/crud-core/repositories"
	models "github.com/yzimhao/trading_engine/v2/internal/models/kline"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type gormKlineRepo struct {
	*repositories.MapperRepository[models.Kline, models.CreateKline, models.UpdateKline, entities.Kline, entities.Kline, map[string]any]
	logger *zap.Logger
}

func NewKlineRepo(datasource datasource.DataSource[gorm.DB], cache cache.Cache, logger *zap.Logger) persistence.KlineRepository {
	cacheWrapperRepo := repositories.NewCacheRepository(
		k_repo.NewGormCrudRepository[entities.Kline, entities.Kline, map[string]any](datasource),
		cache,
	)

	mapperRepo := repositories.NewMapperRepository(
		cacheWrapperRepo,
		b_mappers.NewJSONMapper[models.Kline, models.CreateKline, models.UpdateKline, entities.Kline, entities.Kline, map[string]any](),
	)

	return &gormKlineRepo{
		MapperRepository: mapperRepo,
		logger:           logger,
	}
}
