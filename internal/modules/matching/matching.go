package matching

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/duolacloud/broker-core"
	ds_types "github.com/duolacloud/crud-core/types"
	models_types "github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/pkg/matching"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type inContext struct {
	fx.In
	Broker           broker.Broker
	Logger           *zap.Logger
	TradeVarietyRepo persistence.TradeVarietyRepository
}

type Matching struct {
	broker           broker.Broker
	logger           *zap.Logger
	tradeVarietyRepo persistence.TradeVarietyRepository
	tradePairs       sync.Map
}

func NewMatching(in inContext) *Matching {
	return &Matching{
		broker:           in.Broker,
		logger:           in.Logger,
		tradeVarietyRepo: in.TradeVarietyRepo,
	}
}

func (s *Matching) InitEngine() {
	s.logger.Sugar().Infof("init matching engine")
	// load trade pair
	var (
		cursor string
		next   bool = true
	)
	for next {
		tradeVarieties, extra, err := s.tradeVarietyRepo.CursorQuery(context.Background(), &ds_types.CursorQuery{
			Cursor: cursor,
			Limit:  10,
			Filter: map[string]any{
				"status": map[string]any{
					"eq": models_types.StatusEnabled,
				},
			},
		})

		if err != nil {
			s.logger.Sugar().Errorf("query trade variety error: %v", err)
			continue
		}

		cursor = extra.EndCursor
		next = extra.HasNext

		for _, tradeVariety := range tradeVarieties {
			opts := []matching.Option{
				matching.WithPriceDecimals(int32(tradeVariety.PriceDecimals)),
				matching.WithQuantityDecimals(int32(tradeVariety.QtyDecimals)),
			}
			s.tradePairs.Store(tradeVariety.Symbol, matching.NewEngine(context.Background(), tradeVariety.Symbol, opts...))
			s.logger.Sugar().Infof("init matching engine for symbol: %s", tradeVariety.Symbol)
		}
	}

}

func (s *Matching) Subscribe() {
	s.broker.Subscribe(models_types.TOPIC_ORDER_NEW, s.OnNewOrder)
	s.broker.Subscribe(models_types.TOPIC_NOTIFY_ORDER_CANCEL, s.OnNotifyCancelOrder)
}

func (s *Matching) OnNewOrder(ctx context.Context, event broker.Event) error {
	s.logger.Sugar().Debugf("on new order: %v", event)

	var order models_types.EventOrderNew
	if err := json.Unmarshal(event.Message().Body, &order); err != nil {
		s.logger.Sugar().Errorf("matching new order unmarshal error: %v body: %s", err, string(event.Message().Body))
		return err
	}

	var item matching.QueueItem
	if order.OrderType == matching_types.OrderTypeLimit {
		if order.OrderSide == matching_types.OrderSideSell {
			item = matching.NewAskLimitItem(order.OrderId, *order.Price, *order.Quantity, order.NanoTime)
		} else {
			item = matching.NewBidLimitItem(order.OrderId, *order.Price, *order.Quantity, order.NanoTime)
		}
	} else if order.OrderType == matching_types.OrderTypeMarket {
		// 按成交金额
		if order.Amount != nil {
			if order.OrderSide == matching_types.OrderSideSell {
				item = matching.NewAskMarketAmountItem(order.OrderId, *order.Amount, *order.MaxAmount, order.NanoTime)
			} else {
				item = matching.NewBidMarketAmountItem(order.OrderId, *order.Amount, order.NanoTime)
			}
		} else {
			// 按成交量
			if order.OrderSide == matching_types.OrderSideSell {
				item = matching.NewAskMarketQtyItem(order.OrderId, *order.MaxQty, order.NanoTime)
			} else {
				item = matching.NewBidMarketQtyItem(order.OrderId, *order.Quantity, *order.MaxQty, order.NanoTime)
			}
		}
	}

	if engine := s.engine(order.Symbol); engine != nil {
		engine.AddItem(item)
	}
	return nil
}

func (s *Matching) OnNotifyCancelOrder(ctx context.Context, event broker.Event) error {
	return nil
}

func (s *Matching) engine(symbol string) *matching.Engine {
	if engine, ok := s.tradePairs.Load(symbol); ok {
		return engine.(*matching.Engine)
	}
	return nil
}
