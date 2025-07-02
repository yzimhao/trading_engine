package settlement

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	lockKey = "settle.lock.%s"
)

type SettleLocker struct {
	redis  *redis.Client
	logger *zap.Logger
	mx     sync.Mutex
}

func NewSettleLocker(redis *redis.Client, logger *zap.Logger) *SettleLocker {
	return &SettleLocker{redis: redis, logger: logger}
}

func (s *SettleLocker) Lock(ctx context.Context, orderIds ...string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	for _, id := range orderIds {
		key := fmt.Sprintf(lockKey, id)
		if _, err := s.redis.Do(ctx, "INCR", key).Result(); err != nil {
			s.logger.Sugar().Errorf("settlelock order %s fail err: %s", id, err.Error())
			return err
		}
	}
	return nil
}
func (s *SettleLocker) Unlock(ctx context.Context, orderIds ...string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	for _, id := range orderIds {
		key := fmt.Sprintf(lockKey, id)

		if _, err := s.redis.Do(ctx, "DECR", key).Result(); err != nil {
			s.logger.Sugar().Errorf("settle unlock %s err: %s", id, err.Error())
			return err
		}

		if n, _ := s.redis.Do(ctx, "GET", key).Int64(); n == 0 {
			if _, err := s.redis.Do(ctx, "del", key).Result(); err != nil {
				s.logger.Sugar().Errorf("settle unlock %s fail err: %s", id, err.Error())
				return err
			}
		}
	}
	return nil
}

func (s *SettleLocker) IsExistLock(ctx context.Context, orderId string) (bool, error) {
	key := fmt.Sprintf(lockKey, orderId)
	return s.redis.Do(ctx, "EXISTS", key).Bool()
}
