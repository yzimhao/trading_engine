package gorm

import (
	"context"

	k_repo "github.com/duolacloud/crud-core-gorm/repositories"
	"github.com/duolacloud/crud-core/cache"
	"github.com/duolacloud/crud-core/datasource"
	b_mappers "github.com/duolacloud/crud-core/mappers"
	"github.com/duolacloud/crud-core/repositories"
	"github.com/pkg/errors"
	models "github.com/yzimhao/trading_engine/v2/internal/models/kline"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type gormKlineRepo struct {
	*repositories.MapperRepository[models.Kline, models.CreateKline, models.UpdateKline, entities.Kline, entities.Kline, map[string]any]
	logger     *zap.Logger
	datasource datasource.DataSource[gorm.DB]
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
		datasource:       datasource,
	}
}

func (repo *gormKlineRepo) Save(ctx context.Context, kline *entities.Kline) error {

	db, err := repo.datasource.GetDB(ctx)
	if err != nil {
		return err
	}

	entity := entities.Kline{Period: kline.Period, Symbol: kline.Symbol}

	if !db.Migrator().HasTable(entity.TableName()) {
		if err := db.Table(entity.TableName()).AutoMigrate(&entity); err != nil {
			return errors.Wrap(err, "auto migrate kline table failed")
		}
	}

	var count int64
	query := db.Table(entity.TableName()).
		Where("open_at=? and close_at=?", kline.OpenAt, kline.CloseAt)

	if query.Count(&count); count > 0 {
		if err := db.Table(entity.TableName()).
			Where("open_at=? and close_at=?", kline.OpenAt, kline.CloseAt).
			Updates(map[string]any{
				"open":   kline.Open,
				"high":   kline.High,
				"low":    kline.Low,
				"close":  kline.Close,
				"volume": kline.Volume,
				"amount": kline.Amount,
			}).Error; err != nil {
			return errors.Wrap(err, "update kline failed")
		}
	} else {
		if err := query.Create(kline).Error; err != nil {
			return errors.Wrap(err, "create kline failed")
		}
	}

	return nil
}
