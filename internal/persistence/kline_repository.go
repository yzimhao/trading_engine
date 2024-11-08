package persistence

import (
	"context"

	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
)

type KlineRepository interface {
	// repositories.CrudRepository[models.Kline, models.CreateKline, models.UpdateKline]
	Save(ctx context.Context, kline *entities.Kline) error
}
