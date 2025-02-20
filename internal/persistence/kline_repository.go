package persistence

import (
	"context"

	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	kline_types "github.com/yzimhao/trading_engine/v2/pkg/kline/types"
)

type KlineRepository interface {
	// repositories.CrudRepository[models.Kline, models.CreateKline, models.UpdateKline]
	Save(ctx context.Context, kline *entities.Kline) error
	Find(ctx context.Context, symbol string, period kline_types.PeriodType, start, end int64, limit int) ([]*entities.Kline, error)
}
