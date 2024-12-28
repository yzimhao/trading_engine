package quote

import (
	"context"
	"encoding/json"

	"github.com/duolacloud/broker-core"
	"github.com/redis/go-redis/v9"
	"github.com/yzimhao/trading_engine/v2/app/webws"
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
	ws     *webws.WsManager
}

type inContext struct {
	fx.In
	Logger *zap.Logger
	Broker broker.Broker
	Redis  *redis.Client
	Repo   persistence.KlineRepository
	Ws     *webws.WsManager
}

func NewQuote(in inContext) *Quote {
	return &Quote{
		logger: in.Logger,
		broker: in.Broker,
		redis:  in.Redis,
		repo:   in.Repo,
		ws:     in.Ws,
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

func (q *Quote) processQuote(ctx context.Context, notify models_types.EventNotifyQuote) error {

	q.logger.Sugar().Infof("process quote: %+v", notify)
	k := kline.NewKLine(q.redis, notify.Symbol)

	for _, period := range kline_types.Periods() {
		q.logger.Sugar().Infof("get kline data period: %+v", period)
		// TODO concurrency
		data, err := k.GetData(ctx, period, notify.TradeResult)
		if err != nil {
			q.logger.Sugar().Errorf("get kline data error: %v notifyQuote: %v", err, notify)
			continue
		}

		q.logger.Sugar().Infof("save kline data: %+v", data)

		if err := q.repo.Save(ctx, &entities.Kline{
			OpenAt:  data.OpenAt,
			CloseAt: data.CloseAt,
			Symbol:  notify.Symbol,
			Period:  period,
			Open:    *data.Open,
			High:    *data.High,
			Low:     *data.Low,
			Close:   *data.Close,
			Volume:  *data.Volume,
			Amount:  *data.Amount,
		}); err != nil {
			q.logger.Sugar().Errorf("save kline data error: %v notifyQuote: %v", err, notify)
			continue
		}

		//推送kline记录
		q.ws.Broadcast(ctx, webws.MsgMarketKLineTpl.Format(map[string]string{"period": string(period), "symbol": notify.Symbol}),
			map[string]any{
				"timestamp": data.OpenAt,
				"open":      data.Open,
				"high":      data.High,
				"low":       data.Low,
				"close":     data.Close,
				"volume":    data.Volume,
			},
		)

	}

	return nil
}
