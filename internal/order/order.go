package order

import (
	"context"

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
