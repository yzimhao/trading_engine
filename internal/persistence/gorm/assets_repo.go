package gorm

import (
	k_repo "github.com/duolacloud/crud-core-gorm/repositories"
	"github.com/duolacloud/crud-core/cache"
	"github.com/duolacloud/crud-core/datasource"
	b_mappers "github.com/duolacloud/crud-core/mappers"
	"github.com/duolacloud/crud-core/repositories"
	"github.com/yzimhao/trading_engine/v2/internal/models"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
	"gorm.io/gorm"
)

type gormAssetsRepo struct {
	*repositories.MapperRepository[models.Assets, models.CreateAssets, models.UpdateAssets, entities.Assets, entities.Assets, map[string]any]
}

type gormAssetsLogRepo struct {
	*repositories.MapperRepository[models.AssetsLog, models.CreateAssetsLog, models.UpdateAssetsLog, entities.AssetsLog, entities.AssetsLog, map[string]any]
}

func NewAssetsRepo(datasource datasource.DataSource[gorm.DB], cache cache.Cache) persistence.AssetsRepository {
	cacheWrapperRepo := repositories.NewCacheRepository(
		k_repo.NewGormCrudRepository[entities.Assets, entities.Assets, map[string]any](datasource),
		cache,
	)

	mapperRepo := repositories.NewMapperRepository(
		cacheWrapperRepo,
		b_mappers.NewJSONMapper[models.Assets, models.CreateAssets, models.UpdateAssets, entities.Assets, entities.Assets, map[string]any](),
	)

	return &gormAssetsRepo{
		MapperRepository: mapperRepo,
	}
}

func NewAssetsLogRepo(datasource datasource.DataSource[gorm.DB], cache cache.Cache) *gormAssetsLogRepo {
	cacheWrapperRepo := repositories.NewCacheRepository(
		k_repo.NewGormCrudRepository[entities.AssetsLog, entities.AssetsLog, map[string]any](datasource),
		cache,
	)

	mapperRepo := repositories.NewMapperRepository(
		cacheWrapperRepo,
		b_mappers.NewJSONMapper[models.AssetsLog, models.CreateAssetsLog, models.UpdateAssetsLog, entities.AssetsLog, entities.AssetsLog, map[string]any](),
	)

	return &gormAssetsLogRepo{
		MapperRepository: mapperRepo,
	}
}
