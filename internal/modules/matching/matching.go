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

type Matching struct {
	broker broker.Broker
	logger *zap.Logger
}

func NewMatching(in inContext) *Matching {
	return &Matching{
		broker: in.Broker,
		logger: in.Logger,
	}
}

func (s *Matching) InitEngine() {

}

func (s *Matching) Subscribe() {
	s.broker.Subscribe(types.TOPIC_ORDER_NEW, s.OnNewOrder)
	s.broker.Subscribe(types.TOPIC_ORDER_CANCEL, s.OnCancelOrder)
}

func (s *Matching) OnNewOrder(ctx context.Context, event broker.Event) error {
	return nil
}

func (s *Matching) OnCancelOrder(ctx context.Context, event broker.Event) error {
	return nil
}
