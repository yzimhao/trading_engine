package gorm

import (
	"context"
	"errors"

	k_repo "github.com/duolacloud/crud-core-gorm/repositories"
	"github.com/duolacloud/crud-core/cache"
	"github.com/duolacloud/crud-core/datasource"
	b_mappers "github.com/duolacloud/crud-core/mappers"
	"github.com/duolacloud/crud-core/repositories"
	models "github.com/yzimhao/trading_engine/v2/internal/models/asset"
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type gormAssetRepo struct {
	*repositories.MapperRepository[models.Asset, models.CreateAsset, models.UpdateAsset, entities.Asset, entities.Asset, map[string]any]
	datasource      datasource.DataSource[gorm.DB]
	assetLogRepo    *gormAssetLogRepo
	assetFreezeRepo *gormAssetFreezeRepo
	logger          *zap.Logger
}

type gormAssetLogRepo struct {
	*repositories.MapperRepository[models.AssetLog, models.CreateAssetLog, models.UpdateAssetLog, entities.AssetLog, entities.AssetLog, map[string]any]
}

type gormAssetFreezeRepo struct {
	*repositories.MapperRepository[models.AssetFreeze, models.CreateAssetFreeze, models.UpdateAssetFreeze, entities.AssetFreeze, entities.AssetFreeze, map[string]any]
}

func NewAssetRepo(datasource datasource.DataSource[gorm.DB], cache cache.Cache, logger *zap.Logger) persistence.AssetRepository {
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
		logger:           logger,
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
	db, err := r.datasource.GetDB(ctx)
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		return r.transfer(ctx, tx, symbol, entities.SYSTEM_USER_ID, userId, types.Amount(amount), transId)
	})
}

func (r *gormAssetRepo) Withdraw(ctx context.Context, transId, userId, symbol string, amount string) error {
	db, err := r.datasource.GetDB(ctx)
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		return r.transfer(ctx, tx, symbol, userId, entities.SYSTEM_USER_ID, types.Amount(amount), transId)
	})
}

func (r *gormAssetRepo) Transfer(ctx context.Context, transId, from, to, symbol, amount string) error {
	db, err := r.datasource.GetDB(ctx)
	if err != nil {
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		return r.transfer(ctx, tx, symbol, from, to, types.Amount(amount), transId)
	})

	return err
}

func (r *gormAssetRepo) Freeze(ctx context.Context, transId, userId, symbol string, amount string) error {
	return nil
}

func (r *gormAssetRepo) UnFreeze(ctx context.Context, transId, userId, symbol string, amount string) error {
	return nil
}

// func (r *gormAssetRepo) Query(ctx context.Context, userId, symbol string) (*models.Asset, error) {
// 	return r.QueryOne(ctx, map[string]any{
// 		"symbol": map[string]any{
// 			"eq": symbol,
// 		},
// 		"user_id": map[string]any{
// 			"eq": userId,
// 		},
// 	})
// }

func (r *gormAssetRepo) transfer(ctx context.Context, tx *gorm.DB, symbol, from, to string, amount types.Amount, transId string) error {

	if amount.Cmp(types.Amount("0")) <= 0 {
		return errors.New("amount must be greater than 0")
	}

	//TODO transId去重

	var fromAsset entities.Asset
	//TODO tx.Clauses(clause.Locking{Strength: "FOR UPDATE"})
	if err := tx.First(&fromAsset, "user_id = ? AND symbol = ?", from, symbol).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		} else {
			//根账号自动创建,只有系统账号资产为负数
			if from == entities.SYSTEM_USER_ID {
				fromAsset = entities.Asset{
					UserId: from,
					Symbol: symbol,
				}
			} else {
				return err
			}
		}
	}

	var toAsset entities.Asset
	if err := tx.First(&toAsset, "user_id = ? AND symbol = ?", to, symbol).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		} else {
			toAsset = entities.Asset{
				UserId: to,
				Symbol: symbol,
			}
		}
	}

	fromAsset.TotalBalance = fromAsset.TotalBalance.Sub(amount)
	fromAsset.AvailBalance = fromAsset.AvailBalance.Sub(amount)

	if fromAsset.UserId != entities.SYSTEM_USER_ID {
		if fromAsset.AvailBalance.Cmp(types.Amount("0")) < 0 {
			return errors.New("insufficient balance")
		}
	}

	if tx.Where("user_id = ? AND symbol = ?", from, symbol).FirstOrCreate(&fromAsset).Error != nil {
		return errors.New("update from asset failed")
	}

	toAsset.TotalBalance = toAsset.TotalBalance.Add(amount)
	toAsset.AvailBalance = toAsset.AvailBalance.Add(amount)
	if tx.Where("user_id = ? AND symbol = ?", to, symbol).FirstOrCreate(&toAsset).Error != nil {
		return errors.New("update to asset failed")
	}

	fromLog := &entities.AssetLog{
		UserId:        from,
		Symbol:        symbol,
		BeforeBalance: fromAsset.TotalBalance.Add(amount),
		Amount:        amount.Neg(),
		AfterBalance:  fromAsset.TotalBalance,
		TransID:       transId,
		ChangeType:    entities.AssetChangeTypeTransfer,
	}
	if tx.Create(&fromLog).Error != nil {
		return errors.New("create from asset log failed")
	}

	toLog := &entities.AssetLog{
		UserId:        to,
		Symbol:        symbol,
		BeforeBalance: toAsset.TotalBalance.Sub(amount),
		Amount:        amount,
		AfterBalance:  toAsset.TotalBalance,
		TransID:       transId,
		ChangeType:    entities.AssetChangeTypeTransfer,
	}
	if tx.Create(&toLog).Error != nil {
		return errors.New("create to asset log failed")
	}

	return nil
}
