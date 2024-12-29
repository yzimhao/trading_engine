package gorm

import (
	"context"

	k_repo "github.com/duolacloud/crud-core-gorm/repositories"
	"github.com/duolacloud/crud-core/cache"
	"github.com/duolacloud/crud-core/datasource"
	b_mappers "github.com/duolacloud/crud-core/mappers"
	"github.com/duolacloud/crud-core/repositories"
	"github.com/pkg/errors"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type gormTradeLogRepo struct {
	*repositories.MapperRepository[entities.TradeLog, entities.TradeLog, entities.TradeLog, entities.TradeLog, entities.TradeLog, map[string]any]
	logger     *zap.Logger
	datasource datasource.DataSource[gorm.DB]
}

func NewTradeLogRepo(datasource datasource.DataSource[gorm.DB], cache cache.Cache, logger *zap.Logger) persistence.TradeLogRepository {
	cacheWrapperRepo := repositories.NewCacheRepository(
		k_repo.NewGormCrudRepository[entities.TradeLog, entities.TradeLog, map[string]any](datasource),
		cache,
	)

	mapperRepo := repositories.NewMapperRepository(
		cacheWrapperRepo,
		b_mappers.NewJSONMapper[entities.TradeLog, entities.TradeLog, entities.TradeLog, entities.TradeLog, entities.TradeLog, map[string]any](),
	)

	return &gormTradeLogRepo{
		MapperRepository: mapperRepo,
		logger:           logger,
		datasource:       datasource,
	}
}

func (repo *gormTradeLogRepo) Find(ctx context.Context, symbol string, limit int) ([]*entities.TradeLog, error) {
	//TODO 完善参数功能实现
	db, err := repo.datasource.GetDB(ctx)
	if err != nil {
		return nil, err
	}

	entity := entities.TradeLog{Symbol: symbol}
	if !db.Migrator().HasTable(entity.TableName()) {
		return nil, errors.New("trade log table not found")
	}

	var rows []*entities.TradeLog
	query := db.Table(entity.TableName()).Order("created_at desc").Find(&rows)
	if query.Error != nil {
		return nil, query.Error
	}

	return rows, nil
}
