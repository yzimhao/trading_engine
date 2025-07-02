package settlement

import (
	"context"
	"encoding/json"

	"github.com/duolacloud/broker-core"
	models_types "github.com/yzimhao/trading_engine/v2/internal/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type InContext struct {
	fx.In
	Broker broker.Broker
	Logger *zap.Logger
	Settle *SettleProcessor
}

type SettlementSubscriber struct {
	broker broker.Broker
	logger *zap.Logger
	settle *SettleProcessor
}

func NewSettlementSubscriber(in InContext) *SettlementSubscriber {
	return &SettlementSubscriber{
		broker: in.Broker,
		logger: in.Logger,
		settle: in.Settle,
	}
}

func (s *SettlementSubscriber) Subscribe() {
	s.broker.Subscribe(models_types.TOPIC_ORDER_SETTLE, s.On)
}

func (s *SettlementSubscriber) On(ctx context.Context, event broker.Event) error {
	s.logger.Sugar().Infof("settlement: %+v", event)

	var tradeResult models_types.EventOrderSettle
	if err := json.Unmarshal(event.Message().Body, &tradeResult); err != nil {
		s.logger.Sugar().Errorf("settlement unmarshal error: %v body: %s", err, string(event.Message().Body))
		return err
	}

	return s.settle.Run(ctx, tradeResult.TradeResult)
}
