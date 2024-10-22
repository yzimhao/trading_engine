package order

import (
	"context"

	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderService interface {
	CreateLimit(ctx context.Context, user_id, symbol string, side matching_types.OrderSide, price, qty string) (order *entities.Order, err error)
	CreateMarketByAmount(ctx context.Context, user_id, symbol string, side matching_types.OrderSide, amount string) (order *entities.Order, err error)
	CreateMarketByQty(ctx context.Context, user_id, symbol string, side matching_types.OrderSide, qty string) (order *entities.Order, err error)
	Cancel(ctx context.Context, order_id string, user_id *string) error
}

type orderService struct {
	db     *gorm.DB
	logger *zap.Logger
}

var _ OrderService = (*orderService)(nil)

func NewOrderService(db *gorm.DB, logger *zap.Logger) OrderService {
	return &orderService{
		db:     db,
		logger: logger,
	}
}

func (o *orderService) CreateLimit(ctx context.Context, user_id, symbol string, side matching_types.OrderSide, price, qty string) (order *entities.Order, err error) {
	//TODO implement me
	return nil, nil
}

func (o *orderService) CreateMarketByAmount(ctx context.Context, user_id, symbol string, side matching_types.OrderSide, amount string) (order *entities.Order, err error) {
	//TODO implement me
	return nil, nil
}

func (o *orderService) CreateMarketByQty(ctx context.Context, user_id, symbol string, side matching_types.OrderSide, qty string) (order *entities.Order, err error) {
	//TODO implement me
	return nil, nil
}

func (o *orderService) Cancel(ctx context.Context, order_id string, user_id *string) error {
	//TODO implement me
	return nil
}
