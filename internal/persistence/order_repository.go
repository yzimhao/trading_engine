package persistence

import (
	"context"

	models_types "github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
)

type OrderRepository interface {
	CreateLimit(ctx context.Context, user_id, symbol string, side matching_types.OrderSide, price, qty string) (order *entities.Order, err error)
	CreateMarketByAmount(ctx context.Context, user_id, symbol string, side matching_types.OrderSide, amount string) (order *entities.Order, err error)
	CreateMarketByQty(ctx context.Context, user_id, symbol string, side matching_types.OrderSide, qty string) (order *entities.Order, err error)
	Cancel(ctx context.Context, symbol, order_id string, cancelType models_types.CancelType) error
	LoadUnfinishedOrders(ctx context.Context, symbol string) (orders []*entities.Order, err error)
	HistoryList(ctx context.Context, user_id, symbol string, start, end int64, limit int) (orders []*entities.Order, err error)
}
