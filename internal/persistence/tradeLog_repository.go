package persistence

import (
	"context"

	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
)

type TradeLogRepository interface {
	Find(ctx context.Context, symbol string, limit int) ([]*entities.TradeLog, error)
}
