package order

import (
	"context"

	"github.com/yzimhao/trading_engine/v2/internal/models"
	"go.uber.org/fx"
)

type OrderService interface {
	Create(ctx context.Context) (order_id *string, err error)
}

type InContext struct {
	fx.In
}

type orderService struct{}

var _ OrderService = &orderService{}

func NewOrderService(in InContext) OrderService {
	return &orderService{}
}

func (o *orderService) Create(ctx context.Context) (order_id *string, err error) {
	//TODO implement me
	return nil, nil
}

func (o *orderService) Cancel(ctx context.Context, order_id string, user_id *string) error {

	return nil
}

func (o *orderService) Query(ctx context.Context, order_id string, user_id *string) (*models.Order, error) {

	return nil, nil
}
