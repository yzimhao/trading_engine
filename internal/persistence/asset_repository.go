package persistence

import (
	"context"

	"github.com/duolacloud/crud-core/repositories"
	models "github.com/yzimhao/trading_engine/v2/internal/models/asset"
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
	"gorm.io/gorm"
)

type AssetRepository interface {
	repositories.CrudRepository[models.Asset, models.CreateAsset, models.UpdateAsset]
	Despoit(ctx context.Context, transId, userId, symbol string, amount types.Amount) error
	Withdraw(ctx context.Context, transId, userId, symbol string, amount types.Amount) error
	Transfer(ctx context.Context, transId, from, to, symbol string, amount types.Amount) error
	Freeze(ctx context.Context, tx *gorm.DB, transId, userId, symbol string, amount types.Amount) error
	UnFreeze(ctx context.Context, tx *gorm.DB, transId, userId, symbol string, amount types.Amount) error
}
