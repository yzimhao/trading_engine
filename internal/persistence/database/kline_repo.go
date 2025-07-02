package database

import (
	"context"

	"github.com/pkg/errors"

	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	kline_types "github.com/yzimhao/trading_engine/v2/pkg/kline/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type klineRepo struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewKlineRepo(db *gorm.DB, logger *zap.Logger) persistence.KlineRepository {

	return &klineRepo{
		logger: logger,
		db:     db,
	}
}

func (k *klineRepo) Find(ctx context.Context, symbol string, period kline_types.PeriodType, start, end int64, limit int) ([]*entities.Kline, error) {

	entity := entities.Kline{Symbol: symbol, Period: period}
	if !k.db.Migrator().HasTable(entity.TableName()) {
		return nil, errors.New("kline table not found")
	}

	var rows []*entities.Kline
	query := k.db.Table(entity.TableName()).Order("created_at desc").Find(&rows)
	if query.Error != nil {
		return nil, query.Error
	}

	return rows, nil
}

func (k *klineRepo) Save(ctx context.Context, kline *entities.Kline) error {

	entity := entities.Kline{Period: kline.Period, Symbol: kline.Symbol}

	if !k.db.Migrator().HasTable(entity.TableName()) {
		if err := k.db.Table(entity.TableName()).AutoMigrate(&entity); err != nil {
			return errors.Wrap(err, "auto migrate kline table failed")
		}
	}

	var count int64
	query := k.db.Table(entity.TableName()).
		Where("open_at=? and close_at=?", kline.OpenAt, kline.CloseAt)

	if query.Count(&count); count > 0 {
		if err := k.db.Table(entity.TableName()).
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
