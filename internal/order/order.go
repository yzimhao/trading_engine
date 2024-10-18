package order

import (
	"context"

	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderService interface {
	Create(ctx context.Context) (order_id *string, err error)
}

type InContext struct {
	fx.In
	db     *gorm.DB
	logger *zap.Logger
}

type orderService struct {
	db     *gorm.DB
	logger *zap.Logger
}

var _ OrderService = &orderService{}

func NewOrderService(in InContext) OrderService {
	return &orderService{
		db:     in.db,
		logger: in.logger,
	}
}

func (o *orderService) Create(ctx context.Context) (order_id *string, err error) {
	//TODO implement me
	return nil, nil
}

func (o *orderService) Cancel(ctx context.Context, order_id string, user_id *string) error {

	return nil
}

func (o *orderService) Query(ctx context.Context, order_id string, user_id *string) (*entities.Order, error) {

	return nil, nil
}
