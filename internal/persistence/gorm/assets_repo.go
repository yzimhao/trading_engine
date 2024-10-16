package gorm

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	k_repo "github.com/duolacloud/crud-core-gorm/repositories"
	"github.com/duolacloud/crud-core/cache"
	"github.com/duolacloud/crud-core/datasource"
	b_mappers "github.com/duolacloud/crud-core/mappers"
	"github.com/duolacloud/crud-core/repositories"
	"github.com/yzimhao/trading_engine/v2/internal/models"
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
	"gorm.io/gorm"
)

type gormAssetsRepo struct {
	*repositories.MapperRepository[models.Assets, models.CreateAssets, models.UpdateAssets, entities.Assets, entities.Assets, map[string]any]
	datasource       datasource.DataSource[gorm.DB]
	assetsLogRepo    *gormAssetsLogRepo
	assetsFreezeRepo *gormAssetsFreezeRepo
}

type gormAssetsLogRepo struct {
	*repositories.MapperRepository[models.AssetsLog, models.CreateAssetsLog, models.UpdateAssetsLog, entities.AssetsLog, entities.AssetsLog, map[string]any]
}

type gormAssetsFreezeRepo struct {
	*repositories.MapperRepository[models.AssetsFreeze, models.CreateAssetsFreeze, models.UpdateAssetsFreeze, entities.AssetsFreeze, entities.AssetsFreeze, map[string]any]
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
		datasource:       datasource,
		assetsLogRepo:    newAssetsLogRepo(datasource, cache),
		assetsFreezeRepo: newAssetsFreezeRepo(datasource, cache),
	}

}

func newAssetsLogRepo(datasource datasource.DataSource[gorm.DB], cache cache.Cache) *gormAssetsLogRepo {
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

func newAssetsFreezeRepo(datasource datasource.DataSource[gorm.DB], cache cache.Cache) *gormAssetsFreezeRepo {
	cacheWrapperRepo := repositories.NewCacheRepository(
		k_repo.NewGormCrudRepository[entities.AssetsFreeze, entities.AssetsFreeze, map[string]any](datasource),
		cache,
	)

	mapperRepo := repositories.NewMapperRepository(
		cacheWrapperRepo,
		b_mappers.NewJSONMapper[models.AssetsFreeze, models.CreateAssetsFreeze, models.UpdateAssetsFreeze, entities.AssetsFreeze, entities.AssetsFreeze, map[string]any](),
	)

	return &gormAssetsFreezeRepo{
		MapperRepository: mapperRepo,
	}
}

func (r *gormAssetsRepo) Despoit(ctx context.Context, userId, symbol string, amount string) error {
	return r.transfer(ctx, symbol, entities.SYSTEM_USER_ID, userId, amount, "despoit")

}

func (r *gormAssetsRepo) Withdraw(ctx context.Context, userId, symbol, amount string) error {
	return r.transfer(ctx, symbol, userId, entities.SYSTEM_USER_ID, amount, "withdraw")
}

func (r *gormAssetsRepo) transfer(ctx context.Context, symbol, from, to string, amount string, transId string) error {
	fromUser, err := r.MapperRepository.QueryOne(ctx, map[string]any{
		"and": []map[string]any{
			{
				"user_id": map[string]any{"eq": from},
			},
			{
				"symbol": map[string]any{"eq": symbol},
			},
		},
	})

	if err != nil {
		if err != errors.New("not found") {
			// return errors.Wrap(err, "query from user")
		}
	}

	if fromUser == nil {
		if from != entities.SYSTEM_USER_ID {
			return errors.New("from user not found")
		}

		createFromUser := &models.CreateAssets{
			UserId: from,
			Symbol: symbol,
		}
		fromUser, err = r.MapperRepository.Create(ctx, createFromUser)
		if err != nil {
			return errors.Wrap(err, "create from user")
		}
	}

	toUser, err := r.MapperRepository.QueryOne(ctx, map[string]any{
		"and": []map[string]any{
			{
				"user_id": map[string]any{"eq": to},
			},
			{
				"symbol": map[string]any{"eq": symbol},
			},
		},
	})

	if err != nil {
		if err != errors.New("not found") {
			// return errors.Wrap(err, "query to user")
		}
	}

	if toUser == nil {
		createToUser := &models.CreateAssets{
			UserId: to,
			Symbol: symbol,
		}
		toUser, err = r.MapperRepository.Create(ctx, createToUser)
		if err != nil {
			return errors.Wrap(err, "create to user")
		}
	}

	amountValue := types.Amount(amount)

	fromUserUpdate := &models.UpdateAssets{
		AvailBalance: func() *types.Amount {
			val := fromUser.AvailBalance.Sub(amountValue)
			return &val
		}(),
		TotalBalance: func() *types.Amount {
			val := fromUser.TotalBalance.Sub(amountValue)
			return &val
		}(),
	}
	toUserUpdate := &models.UpdateAssets{
		AvailBalance: func() *types.Amount {
			val := toUser.AvailBalance.Add(amountValue)
			return &val
		}(),
		TotalBalance: func() *types.Amount {
			val := toUser.TotalBalance.Add(amountValue)
			return &val
		}(),
	}

	//TODO 校验fromUser.Available >= 0

	//logs
	fromLog := models.CreateAssetsLog{
		UserId:  fromUser.UserId,
		Symbol:  fromUser.Symbol,
		Amount:  fmt.Sprintf("-%s", amountValue),
		TransID: transId,
		// ChangeType:    ,
	}
	toLog := models.CreateAssetsLog{
		UserId:  toUser.UserId,
		Symbol:  toUser.Symbol,
		Amount:  amountValue.String(),
		TransID: transId,
		// ChangeType:    ,
	}

	_, err = r.assetsLogRepo.Create(ctx, &fromLog)
	if err != nil {
		return errors.Wrap(err, "create from log")
	}

	_, err = r.assetsLogRepo.Create(ctx, &toLog)
	if err != nil {
		return errors.Wrap(err, "create to log")
	}

	_, err = r.MapperRepository.Update(ctx, fromUser.Id, fromUserUpdate)
	if err != nil {
		return errors.Wrap(err, "update from user")
	}

	_, err = r.MapperRepository.Update(ctx, toUser.Id, toUserUpdate)
	if err != nil {
		return errors.Wrap(err, "update to user")
	}

	return nil
}
