package persistence

import (
	"context"

	"github.com/duolacloud/crud-core/repositories"
	models "github.com/yzimhao/trading_engine/v2/internal/models/asset"
)

type AssetRepository interface {
	repositories.CrudRepository[models.Asset, models.CreateAsset, models.UpdateAsset]
	Despoit(ctx context.Context, transId, userId, symbol, amount string) error
	Withdraw(ctx context.Context, transId, userId, symbol, amount string) error
	Transfer(ctx context.Context, transId, from, to, symbol, amount string) error
	Freeze(ctx context.Context, transId, userId, symbol, amount string) error
	UnFreeze(ctx context.Context, transId, userId, symbol, amount string) error
}
