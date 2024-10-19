package matching

import (
	"context"

	"github.com/duolacloud/broker-core"
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type inContext struct {
	fx.In
	Broker broker.Broker
	Logger *zap.Logger
}

type MatchingSubscriber struct {
	broker broker.Broker
	logger *zap.Logger
}

func NewMatchingSubscriber(in inContext) *MatchingSubscriber {
	return &MatchingSubscriber{
		broker: in.Broker,
		logger: in.Logger,
	}
}

func (s *MatchingSubscriber) Subscribe() {
	s.broker.Subscribe(types.TOPIC_ORDER_NEW, s.On)
}

func (s *MatchingSubscriber) On(ctx context.Context, event broker.Event) error {
	return nil
}
