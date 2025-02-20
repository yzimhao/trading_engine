package database

import (
	"context"
	"errors"

	k_repo "github.com/duolacloud/crud-core-gorm/repositories"
	"github.com/duolacloud/crud-core/cache"
	"github.com/duolacloud/crud-core/datasource"
	b_mappers "github.com/duolacloud/crud-core/mappers"
	"github.com/duolacloud/crud-core/repositories"
	datasource_types "github.com/duolacloud/crud-core/types"
	models "github.com/yzimhao/trading_engine/v2/internal/models/asset"
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *gormAssetRepo) Despoit(ctx context.Context, transId, userId, symbol string, amount types.Numeric) error {
	db, err := r.datasource.GetDB(ctx)
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		return r.transfer(ctx, tx, symbol, entities.SYSTEM_USER_ROOT, userId, amount, transId)
	})
}

func (r *gormAssetRepo) Withdraw(ctx context.Context, transId, userId, symbol string, amount types.Numeric) error {
	db, err := r.datasource.GetDB(ctx)
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		return r.transfer(ctx, tx, symbol, userId, entities.SYSTEM_USER_ROOT, amount, transId)
	})
}

// 两个user之间的转账
func (r *gormAssetRepo) Transfer(ctx context.Context, transId, from, to, symbol string, amount types.Numeric) error {
	db, err := r.datasource.GetDB(ctx)
	if err != nil {
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		return r.transfer(ctx, tx, symbol, from, to, amount, transId)
	})

	return err
}

// 冻结资产
// 这里使用tx传入，方便在结算的时候事务中使用
func (r *gormAssetRepo) Freeze(ctx context.Context, tx *gorm.DB, transId, userId, symbol string, amount types.Numeric) (*entities.AssetFreeze, error) {
	if amount.Cmp(types.NumericZero) < 0 {
		return nil, errors.New("amount must be >= 0")
	}

	asset := entities.Asset{UserId: userId, Symbol: symbol}
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ? AND symbol = ?", userId, symbol).FirstOrCreate(&asset).Error; err != nil {
		return nil, err
	}

	//冻结金额为0，冻结全部可用
	if amount.Cmp(types.NumericZero) == 0 {
		amount = asset.AvailBalance
	}

	asset.AvailBalance = asset.AvailBalance.Sub(amount)
	asset.FreezeBalance = asset.FreezeBalance.Add(amount)

	if asset.AvailBalance.Cmp(types.NumericZero) < 0 {
		return nil, errors.New("insufficient balance")
	}

	if tx.Where("user_id = ? AND symbol = ?", userId, symbol).Updates(&asset).Error != nil {
		return nil, errors.New("update asset failed")
	}

	//freeze log
	freezeLog := &entities.AssetFreeze{
		UserId:       userId,
		Symbol:       symbol,
		Amount:       amount,
		FreezeAmount: amount,
		TransId:      transId,
		// FreezeType:   entities.FreezeTypeTrade, //TODO 冻结类型
	}
	if tx.Create(&freezeLog).Error != nil {
		return nil, errors.New("create freeze log failed")
	}

	return freezeLog, nil
}

// 解冻资产
// amount为0，则解冻这条记录的全部剩余
func (r *gormAssetRepo) UnFreeze(ctx context.Context, tx *gorm.DB, transId, userId, symbol string, amount types.Numeric) error {
	if amount.Cmp(types.NumericZero) < 0 {
		return errors.New("amount must be > 0")
	}

	freeze := entities.AssetFreeze{UserId: userId, Symbol: symbol, TransId: transId}
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ? AND symbol = ? AND trans_id = ?", userId, symbol, transId).First(&freeze).Error; err != nil {
		return err
	}

	if freeze.Status == entities.FreezeStatusDone {
		return errors.New("unfreeze already done")
	}

	//解冻金额为0，则全部金额解冻
	if amount.Cmp(types.NumericZero) == 0 {
		amount = freeze.FreezeAmount
	}

	freeze.FreezeAmount = freeze.FreezeAmount.Sub(amount)
	if freeze.FreezeAmount.Cmp(types.NumericZero) == 0 {
		freeze.Status = entities.FreezeStatusDone
	}

	if tx.Where("user_id = ? AND symbol = ? AND trans_id = ?", userId, symbol, transId).Updates(&freeze).Error != nil {
		return errors.New("update freeze failed")
	}
	// 冻结资产变可用资产
	asset := entities.Asset{UserId: userId, Symbol: symbol}
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ? AND symbol = ?", userId, symbol).FirstOrCreate(&asset).Error; err != nil {
		return err
	}
	asset.FreezeBalance = asset.FreezeBalance.Sub(amount)
	asset.AvailBalance = asset.AvailBalance.Add(amount)
	if tx.Where("user_id = ? AND symbol = ?", userId, symbol).Updates(&asset).Error != nil {
		return errors.New("update asset failed")
	}

	return nil
}

func (r *gormAssetRepo) QueryFreeze(ctx context.Context, filter map[string]any) (assetFreezes []*models.AssetFreeze, err error) {
	query := &datasource_types.PageQuery{
		Filter: filter,
	}
	data, err := r.assetFreezeRepo.Query(ctx, query)
	return data, err
}

func (r *gormAssetRepo) TransferWithTx(ctx context.Context, tx *gorm.DB, transId, from, to, symbol string, amount types.Numeric) error {
	return r.transfer(ctx, tx, symbol, from, to, amount, transId)
}

func (r *gormAssetRepo) transfer(ctx context.Context, tx *gorm.DB, symbol, from, to string, amount types.Numeric, transId string) error {
	if from == to {
		return errors.New("from and to cannot be the same")
	}

	if amount.Cmp(types.NumericZero) <= 0 {
		return errors.New("amount must be greater than 0")
	}

	//TODO transId去重

	fromAsset := entities.Asset{UserId: from, Symbol: symbol}
	//TODO tx.Clauses(clause.Locking{Strength: "FOR UPDATE"})
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ? AND symbol = ?", from, symbol).FirstOrCreate(&fromAsset).Error; err != nil {
		return err
	}

	toAsset := entities.Asset{UserId: to, Symbol: symbol}
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ? AND symbol = ?", to, symbol).FirstOrCreate(&toAsset).Error; err != nil {
		return err
	}

	fromAsset.TotalBalance = fromAsset.TotalBalance.Sub(amount)
	fromAsset.AvailBalance = fromAsset.AvailBalance.Sub(amount)

	if fromAsset.UserId != entities.SYSTEM_USER_ROOT {
		if fromAsset.AvailBalance.Cmp(types.NumericZero) < 0 {
			return errors.New("insufficient balance")
		}
	}

	if tx.Where("user_id = ? AND symbol = ?", from, symbol).Updates(&fromAsset).Error != nil {
		return errors.New("update from asset failed")
	}

	toAsset.TotalBalance = toAsset.TotalBalance.Add(amount)
	toAsset.AvailBalance = toAsset.AvailBalance.Add(amount)
	if tx.Where("user_id = ? AND symbol = ?", to, symbol).Updates(&toAsset).Error != nil {
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
