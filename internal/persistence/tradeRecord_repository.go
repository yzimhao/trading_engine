package persistence

import (
	"context"

	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
)

type TradeRecordRepository interface {
	Find(ctx context.Context, symbol string, limit int) ([]*entities.TradeRecord, error)
}
