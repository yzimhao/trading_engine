package gorm

import (
	"context"

	k_repo "github.com/duolacloud/crud-core-gorm/repositories"
	"github.com/duolacloud/crud-core/cache"
	"github.com/duolacloud/crud-core/datasource"
	b_mappers "github.com/duolacloud/crud-core/mappers"
	"github.com/duolacloud/crud-core/repositories"
	models "github.com/yzimhao/trading_engine/v2/internal/models/asset"
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
	"gorm.io/gorm"
)

type gormAssetRepo struct {
	*repositories.MapperRepository[models.Asset, models.CreateAsset, models.UpdateAsset, entities.Asset, entities.Asset, map[string]any]
	datasource      datasource.DataSource[gorm.DB]
	assetLogRepo    *gormAssetLogRepo
	assetFreezeRepo *gormAssetFreezeRepo
}

type gormAssetLogRepo struct {
	*repositories.MapperRepository[models.AssetLog, models.CreateAssetLog, models.UpdateAssetLog, entities.AssetLog, entities.AssetLog, map[string]any]
}

type gormAssetFreezeRepo struct {
	*repositories.MapperRepository[models.AssetFreeze, models.CreateAssetFreeze, models.UpdateAssetFreeze, entities.AssetFreeze, entities.AssetFreeze, map[string]any]
}

func NewAssetRepo(datasource datasource.DataSource[gorm.DB], cache cache.Cache) persistence.AssetRepository {
	cacheWrapperRepo := repositories.NewCacheRepository(
		k_repo.NewGormCrudRepository[entities.Asset, entities.Asset, map[string]any](datasource),
		cache,
	)

	mapperRepo := repositories.NewMapperRepository(
		cacheWrapperRepo,
		b_mappers.NewJSONMapper[models.Asset, models.CreateAsset, models.UpdateAsset, entities.Asset, entities.Asset, map[string]any](),
	)

	return &gormAssetRepo{
		MapperRepository: mapperRepo,
		datasource:       datasource,
		assetLogRepo:     newAssetLogRepo(datasource, cache),
		assetFreezeRepo:  newAssetFreezeRepo(datasource, cache),
	}

}

func newAssetLogRepo(datasource datasource.DataSource[gorm.DB], cache cache.Cache) *gormAssetLogRepo {
	cacheWrapperRepo := repositories.NewCacheRepository(
		k_repo.NewGormCrudRepository[entities.AssetLog, entities.AssetLog, map[string]any](datasource),
		cache,
	)

	mapperRepo := repositories.NewMapperRepository(
		cacheWrapperRepo,
		b_mappers.NewJSONMapper[models.AssetLog, models.CreateAssetLog, models.UpdateAssetLog, entities.AssetLog, entities.AssetLog, map[string]any](),
	)

	return &gormAssetLogRepo{
		MapperRepository: mapperRepo,
	}
}

func newAssetFreezeRepo(datasource datasource.DataSource[gorm.DB], cache cache.Cache) *gormAssetFreezeRepo {
	cacheWrapperRepo := repositories.NewCacheRepository(
		k_repo.NewGormCrudRepository[entities.AssetFreeze, entities.AssetFreeze, map[string]any](datasource),
		cache,
	)

	mapperRepo := repositories.NewMapperRepository(
		cacheWrapperRepo,
		b_mappers.NewJSONMapper[models.AssetFreeze, models.CreateAssetFreeze, models.UpdateAssetFreeze, entities.AssetFreeze, entities.AssetFreeze, map[string]any](),
	)

	return &gormAssetFreezeRepo{
		MapperRepository: mapperRepo,
	}
}

func (r *gormAssetRepo) Despoit(ctx context.Context, transId, userId, symbol string, amount string) error {
	return r.transfer(ctx, symbol, entities.SYSTEM_USER_ID, userId, types.Amount(amount), transId)
}

func (r *gormAssetRepo) Withdraw(ctx context.Context, transId, userId, symbol string, amount string) error {
	return r.transfer(ctx, symbol, userId, entities.SYSTEM_USER_ID, types.Amount(amount), transId)
}

func (r *gormAssetRepo) Transfer(ctx context.Context, transId, from, to, symbol, amount string) error {
	return r.transfer(ctx, symbol, from, to, types.Amount(amount), transId)
}

func (r *gormAssetRepo) Freeze(ctx context.Context, transId, userId, symbol string, amount string) error {
	return nil
}

func (r *gormAssetRepo) UnFreeze(ctx context.Context, transId, userId, symbol string, amount string) error {
	return nil
}

func (r *gormAssetRepo) transfer(ctx context.Context, symbol, from, to string, amount types.Amount, transId string) error {
	return nil
}
