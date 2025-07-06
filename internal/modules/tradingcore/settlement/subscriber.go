package settlement

import (
	"context"
	"encoding/json"

	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	models_types "github.com/yzimhao/trading_engine/v2/internal/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type InContext struct {
	fx.In
	Consume *provider.Consume
	Logger  *zap.Logger
	Settle  *SettleProcessor
}

type SettlementSubscriber struct {
	consume *provider.Consume
	logger  *zap.Logger
	settle  *SettleProcessor
}

func NewSettlementSubscriber(in InContext) *SettlementSubscriber {
	return &SettlementSubscriber{
		consume: in.Consume,
		logger:  in.Logger,
		settle:  in.Settle,
	}
}

func (s *SettlementSubscriber) Subscribe() {
	// s.broker.Subscribe(models_types.TOPIC_ORDER_SETTLE, s.process)
	s.consume.Subscribe(models_types.TOPIC_ORDER_SETTLE, func(ctx context.Context, data []byte) {
		if err := s.process(ctx, data); err != nil {
			s.logger.Sugar().Errorf("settlement process: %s err: %s", data, err)
		}
	})
}

func (s *SettlementSubscriber) process(ctx context.Context, msg []byte) error {
	s.logger.Sugar().Infof("settlement: %s", msg)

	var tradeResult models_types.EventOrderSettle
	if err := json.Unmarshal(msg, &tradeResult); err != nil {
		s.logger.Sugar().Errorf("settlement unmarshal error: %v body: %s", err, msg)
		return err
	}

	return s.settle.Run(ctx, tradeResult.TradeResult)
}
