package persistence

import (
	"context"

	"github.com/duolacloud/crud-core/repositories"
	"github.com/yzimhao/trading_engine/v2/internal/models"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
)

type AssetsRepository interface {
	repositories.CrudRepository[models.Assets, models.CreateAssets, models.UpdateAssets]
	Despoit(ctx context.Context, userId, symbol, amount string) error
	Withdraw(ctx context.Context, userId, symbol, amount string) error
	FindOne(ctx context.Context, userId, symbol string) (*entities.Assets, error)
	Transfer(ctx context.Context, from, to, symbol, amount string) error
	FindAssetHistory(ctx context.Context) ([]entities.AssetsLog, error)
}
