package provider

//暂时先用redis的list，后期优化rocketmq的参数，全部替换
import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Produce struct {
	redis  *redis.Client
	logger *zap.Logger
}
type Consume struct {
	redis  *redis.Client
	logger *zap.Logger
}

func NewProduce(r *redis.Client, l *zap.Logger) *Produce {
	return &Produce{
		redis:  r,
		logger: l,
	}
}
func NewConsume(r *redis.Client, l *zap.Logger) *Consume {
	return &Consume{
		redis:  r,
		logger: l,
	}
}

func (p *Produce) Publish(ctx context.Context, topic string, data any) error {
	err := p.redis.Do(ctx, "LPUSH", topic, data)
	if topic != "websocket_msg" {
		p.logger.Sugar().Debugf("publish %s msg: %s", topic, data)
	}
	return err.Err()
}

func (c *Consume) Subscribe(topic string, cb func(ctx context.Context, data []byte)) {
	go func() {
		for {
			ctx := context.Background()
			reply, err := c.redis.BRPop(ctx, 0, topic).Result()
			if err != nil {
				continue
			}

			if topic != "websocket_msg" {
				c.logger.Sugar().Debugf("subscribe %s msg: %+v", topic, reply)
			}

			if cb != nil {
				cb(ctx, []byte(reply[1]))
			}
		}
	}()
}
