package settlement

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/modules/tradingcore/orderlock"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type inCancelOrderContext struct {
	fx.In
	Consume   *provider.Consume
	Logger    *zap.Logger
	OrderRepo persistence.OrderRepository
	Redis     *redis.Client
	Locker    *orderlock.OrderLock
}

type CancelOrderSubscriber struct {
	consume   *provider.Consume
	logger    *zap.Logger
	orderRepo persistence.OrderRepository
	redis     *redis.Client
	locker    *orderlock.OrderLock
	maxRetry  int
}

func NewCancelOrderSubscriber(in inCancelOrderContext) *CancelOrderSubscriber {
	return &CancelOrderSubscriber{
		consume:   in.Consume,
		logger:    in.Logger,
		orderRepo: in.OrderRepo,
		redis:     in.Redis,
		locker:    in.Locker,
		maxRetry:  20,
	}
}

func (s *CancelOrderSubscriber) Subscribe() {
	s.consume.Subscribe(types.TOPIC_PROCESS_ORDER_CANCEL, func(ctx context.Context, data []byte) {
		s.On(ctx, data)
	})
}

func (s *CancelOrderSubscriber) On(ctx context.Context, msg []byte) error {
	s.logger.Sugar().Debugf("cancel order %s", msg)

	var data types.EventCancelOrder
	if err := json.Unmarshal(msg, &data); err != nil {
		s.logger.Sugar().Errorf("unmarshal cancel order event error: %v, event: %s", err, msg)
		return err
	}
	return s.process(ctx, data, 0)
}

func (s *CancelOrderSubscriber) process(ctx context.Context, data types.EventCancelOrder, retryCount int) error {
	s.logger.Sugar().Infof("order cancel %s, retry count: %d", data.OrderId, retryCount)
	//锁等待结算那边全部结束才能取消
	ok, err := s.locker.IsLocked(ctx, data.OrderId)
	if err != nil {
		return err
	}

	if ok {
		s.logger.Sugar().Debugf("order cancel %s is locked, retry count: %d", data.OrderId, retryCount)

		if retryCount <= s.maxRetry {
			time.Sleep(time.Duration(500) * time.Millisecond)
			return s.process(ctx, data, retryCount+1)
		}
		s.logger.Sugar().Errorf("order cancel %s is locked, retry over max retry", data.OrderId)
		return errors.New("retry over max retry")
	}

	return s.orderRepo.Cancel(ctx, data.Symbol, data.OrderId, types.CancelTypeUser)
}
