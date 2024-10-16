package persistence

import (
	"context"

	"github.com/duolacloud/crud-core/repositories"
	"github.com/yzimhao/trading_engine/v2/internal/models"
)

type AssetsRepository interface {
	repositories.CrudRepository[models.Assets, models.CreateAssets, models.UpdateAssets]
	Despoit(ctx context.Context, userId, symbol, amount string) error
	Withdraw(ctx context.Context, userId, symbol, amount string) error
}
