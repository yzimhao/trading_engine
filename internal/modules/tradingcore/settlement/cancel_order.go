package settlement

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/duolacloud/broker-core"
	"github.com/redis/go-redis/v9"
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type inCancelOrderContext struct {
	fx.In
	Broker    broker.Broker
	Logger    *zap.Logger
	OrderRepo persistence.OrderRepository
	Redis     *redis.Client
	Locker    *SettleLocker
}

type CancelOrderSubscriber struct {
	broker    broker.Broker
	logger    *zap.Logger
	orderRepo persistence.OrderRepository
	redis     *redis.Client
	locker    *SettleLocker
	maxRetry  int
}

func NewCancelOrderSubscriber(in inCancelOrderContext) *CancelOrderSubscriber {
	return &CancelOrderSubscriber{
		broker:    in.Broker,
		logger:    in.Logger,
		orderRepo: in.OrderRepo,
		redis:     in.Redis,
		locker:    in.Locker,
		maxRetry:  20,
	}
}

func (s *CancelOrderSubscriber) Subscribe() {
	s.broker.Subscribe(types.TOPIC_PROCESS_ORDER_CANCEL, s.On)
}

func (s *CancelOrderSubscriber) On(ctx context.Context, event broker.Event) error {
	s.logger.Info("cancel order", zap.Any("event", event))

	var data types.EventCancelOrder
	if err := json.Unmarshal(event.Message().Body, &data); err != nil {
		s.logger.Sugar().Errorf("unmarshal cancel order event error: %v, event: %v", err, event)
		return err
	}
	return s.process(ctx, data, 0)
}

func (s *CancelOrderSubscriber) process(ctx context.Context, data types.EventCancelOrder, retryCount int) error {
	s.logger.Sugar().Infof("order cancel %s, retry count: %d", data.OrderId, retryCount)
	//锁等待结算那边全部结束才能取消
	ok, err := s.locker.IsExistLock(ctx, data.OrderId)
	if err != nil {
		return err
	}

	if ok {
		s.logger.Sugar().Errorf("order cancel %s is locked, retry count: %d", data.OrderId, retryCount)

		if retryCount <= s.maxRetry {
			time.Sleep(time.Duration(500) * time.Millisecond)
			return s.process(ctx, data, retryCount+1)
		}
		s.logger.Sugar().Errorf("order cancel %s is locked, retry over max retry", data.OrderId)
		return errors.New("retry over max retry")
	}

	return s.orderRepo.Cancel(ctx, data.Symbol, data.OrderId, types.CancelTypeUser)
}
