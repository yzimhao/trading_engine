package quote

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	notification_ws "github.com/yzimhao/trading_engine/v2/internal/modules/notification/ws"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	models_types "github.com/yzimhao/trading_engine/v2/internal/types"
	"github.com/yzimhao/trading_engine/v2/pkg/kline"
	kline_types "github.com/yzimhao/trading_engine/v2/pkg/kline/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Quote struct {
	logger      *zap.Logger
	consume     *provider.Consume
	redis       *redis.Client
	repo        persistence.KlineRepository
	ws          *notification_ws.WsManager
	productRepo persistence.ProductRepository
}

type inContext struct {
	fx.In
	Logger      *zap.Logger
	Consume     *provider.Consume
	Redis       *redis.Client
	Repo        persistence.KlineRepository
	Ws          *notification_ws.WsManager
	ProductRepo persistence.ProductRepository
}

func NewQuote(in inContext) *Quote {
	return &Quote{
		logger:      in.Logger,
		consume:     in.Consume,
		redis:       in.Redis,
		repo:        in.Repo,
		ws:          in.Ws,
		productRepo: in.ProductRepo,
	}
}

func (q *Quote) Subscribe() {
	// q.broker.Subscribe(models_types.TOPIC_NOTIFY_QUOTE, q.OnNotifyQuote)
	q.consume.Subscribe(models_types.TOPIC_NOTIFY_QUOTE, func(ctx context.Context, msg []byte) {
		if err := q.OnNotifyQuote(ctx, msg); err != nil {
			q.logger.Sugar().Errorf("quote subscribe msg: %s err: %s", msg, err)
		}
	})
}

func (q *Quote) OnNotifyQuote(ctx context.Context, msg []byte) error {
	q.logger.Sugar().Debugf("on notify quote: %v", msg)
	var notifyQuote models_types.EventNotifyQuote
	if err := json.Unmarshal(msg, &notifyQuote); err != nil {
		q.logger.Sugar().Errorf("unmarshal notify quote error: %v", err)
		return err
	}

	return q.processQuote(ctx, notifyQuote)
}

func (q *Quote) processQuote(ctx context.Context, notify models_types.EventNotifyQuote) error {

	q.logger.Sugar().Infof("process quote: %+v", notify)
	product, err := q.productRepo.Get(notify.Symbol)
	if err != nil {
		q.logger.Sugar().Errorf("quote process tradelog product.Get error: %v", err)
		return err
	}

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
		q.ws.Broadcast(ctx, notification_ws.MsgMarketKLineTpl.Format(map[string]string{"period": string(period), "symbol": notify.Symbol}),
			[6]any{
				data.OpenAt.UnixMilli(),
				// common.FormatStrNumber(*data.Open, product.PriceDecimals),
				// common.FormatStrNumber(*data.High, product.PriceDecimals),
				// common.FormatStrNumber(*data.Low, product.PriceDecimals),
				// common.FormatStrNumber(*data.Close, product.PriceDecimals),
				// common.FormatStrNumber(*data.Volume, product.QtyDecimals),
				data.Open.Truncate(product.PriceDecimals),
				data.High.Truncate(product.PriceDecimals),
				data.Low.Truncate(product.PriceDecimals),
				data.Close.Truncate(product.PriceDecimals),
				data.Volume.Truncate(product.QtyDecimals),
			},
		)

	}

	return nil
}
