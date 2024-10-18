package settlement

import (
	"context"

	"github.com/duolacloud/broker-core"
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type InContext struct {
	fx.In
	Broker broker.Broker
	Logger *zap.Logger
}

type SettlementSubscriber struct {
	broker broker.Broker
	logger *zap.Logger
}

func NewSettlementSubscriber(in InContext) *SettlementSubscriber {
	return &SettlementSubscriber{
		broker: in.Broker,
		logger: in.Logger,
	}
}

func (s *SettlementSubscriber) Subscribe() {
	s.broker.Subscribe(types.TOPIC_ORDER_SETTLE, s.On)
}

func (s *SettlementSubscriber) On(ctx context.Context, event broker.Event) error {
	s.logger.Info("settlement", zap.Any("event", event))
	//TODO

	return nil
}
