package persistence

import (
	"context"

	"github.com/duolacloud/crud-core/repositories"
	models "github.com/yzimhao/trading_engine/v2/internal/models/asset"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
)

type AssetRepository interface {
	repositories.CrudRepository[models.Asset, models.CreateAsset, models.UpdateAsset]
	Despoit(ctx context.Context, userId, symbol, amount string) (order_id string, err error)
	Withdraw(ctx context.Context, userId, symbol, amount string) (order_id string, err error)
	FindOne(ctx context.Context, userId, symbol string) (*entities.Asset, error)
	Transfer(ctx context.Context, from, to, symbol, amount string) error
	FindAssetHistory(ctx context.Context) ([]entities.AssetLog, error)
}
