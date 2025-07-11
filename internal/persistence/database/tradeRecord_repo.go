package database

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type tradeRecordRepo struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewTradeRecordRepo(datasource *gorm.DB, logger *zap.Logger) persistence.TradeRecordRepository {

	return &tradeRecordRepo{
		logger: logger,
		db:     datasource,
	}
}

func (t *tradeRecordRepo) Find(ctx context.Context, symbol string, limit int) ([]*entities.TradeRecord, error) {

	entity := entities.TradeRecord{Symbol: symbol}
	if !t.db.Migrator().HasTable(entity.TableName()) {
		return nil, errors.New("trade record table not found")
	}

	var rows []*entities.TradeRecord
	query := t.db.Table(entity.TableName()).Limit(limit).Order("created_at desc").Find(&rows)
	if query.Error != nil {
		return nil, query.Error
	}

	return rows, nil
}
