package quote

import (
	"context"
	"encoding/json"

	"github.com/duolacloud/broker-core"
	"github.com/redis/go-redis/v9"
	models_types "github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
	"github.com/yzimhao/trading_engine/v2/pkg/kline"
	kline_types "github.com/yzimhao/trading_engine/v2/pkg/kline/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Quote struct {
	logger *zap.Logger
	broker broker.Broker
	redis  *redis.Client
	repo   persistence.KlineRepository
}

type inContext struct {
	fx.In
	Logger *zap.Logger
	Broker broker.Broker
	Redis  *redis.Client
	Repo   persistence.KlineRepository
}

func NewQuote(in inContext) *Quote {
	return &Quote{
		logger: in.Logger,
		broker: in.Broker,
		redis:  in.Redis,
		repo:   in.Repo,
	}
}

func (q *Quote) Subscribe() {
	q.broker.Subscribe(models_types.TOPIC_NOTIFY_QUOTE, q.OnNotifyQuote)
}

func (q *Quote) OnNotifyQuote(ctx context.Context, event broker.Event) error {
	q.logger.Sugar().Debugf("on notify quote: %v", event)
	var notifyQuote models_types.EventNotifyQuote
	if err := json.Unmarshal(event.Message().Body, &notifyQuote); err != nil {
		q.logger.Sugar().Errorf("unmarshal notify quote error: %v", err)
		return err
	}

	return q.processQuote(ctx, notifyQuote)
}

func (q *Quote) processQuote(ctx context.Context, notifyQuote models_types.EventNotifyQuote) error {

	k := kline.NewKLine(q.redis, notifyQuote.Symbol)
	for _, period := range kline_types.Periods() {

		data, err := k.GetData(ctx, period, notifyQuote.TradeResult)
		if err != nil {
			q.logger.Sugar().Errorf("get kline data error: %v notifyQuote: %v", err, notifyQuote)
			continue
		}

		if err := q.repo.Save(ctx, &entities.Kline{
			Symbol: notifyQuote.Symbol,
			Period: period,
			Open:   *data.Open,
			High:   *data.High,
			Low:    *data.Low,
			Close:  *data.Close,
			Volume: *data.Volume,
			Amount: *data.Amount,
		}); err != nil {
			q.logger.Sugar().Errorf("save kline data error: %v notifyQuote: %v", err, notifyQuote)
			continue
		}
	}

	return nil
}
