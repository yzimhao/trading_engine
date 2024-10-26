package settlement

import (
	"context"

	"github.com/duolacloud/broker-core"
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type inCancelOrderContext struct {
	fx.In
	Broker broker.Broker
	Logger *zap.Logger
}

type CancelOrderSubscriber struct {
	broker broker.Broker
	logger *zap.Logger
}

func NewCancelOrderSubscriber(in inCancelOrderContext) *CancelOrderSubscriber {
	return &CancelOrderSubscriber{
		broker: in.Broker,
		logger: in.Logger,
	}
}

func (s *CancelOrderSubscriber) Subscribe() {
	s.broker.Subscribe(types.TOPIC_PROCESS_ORDER_CANCEL, s.On)
}

func (s *CancelOrderSubscriber) On(ctx context.Context, event broker.Event) error {
	s.logger.Info("cancel order", zap.Any("event", event))
	//TODO

	return nil
}
