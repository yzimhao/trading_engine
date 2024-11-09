package matching

import (
	"context"
	"encoding/json"
	"slices"
	"sync"

	"github.com/duolacloud/broker-core"
	ds_types "github.com/duolacloud/crud-core/types"
	"github.com/spf13/viper"
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
	Viper            *viper.Viper
}

type Matching struct {
	broker           broker.Broker
	logger           *zap.Logger
	tradeVarietyRepo persistence.TradeVarietyRepository
	tradePairs       sync.Map
	viper            *viper.Viper
}

func NewMatching(in inContext) *Matching {
	return &Matching{
		broker:           in.Broker,
		logger:           in.Logger,
		tradeVarietyRepo: in.TradeVarietyRepo,
		viper:            in.Viper,
	}
}

func (s *Matching) InitEngine() {
	s.logger.Sugar().Infof("init matching engine")
	localSymbols := s.viper.GetStringSlice("matching.local_symbols")

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
			if len(localSymbols) > 0 {
				if !slices.Contains(localSymbols, tradeVariety.Symbol) {
					continue
				}
			}

			opts := []matching.Option{
				matching.WithPriceDecimals(int32(tradeVariety.PriceDecimals)),
				matching.WithQuantityDecimals(int32(tradeVariety.QtyDecimals)),
				// matching.WithLogger(s.logger),
			}
			engine := matching.NewEngine(context.Background(), tradeVariety.Symbol, opts...)

			engine.OnRemoveResult(func(result matching_types.RemoveResult) {
				s.logger.Sugar().Infof("symbol: %s remove result: %v", result.Symbol, result)
				s.processCancelOrderResult(result)
			})
			engine.OnTradeResult(func(result matching_types.TradeResult) {
				s.logger.Sugar().Infof("symbol: %s trade result: %v", result.Symbol, result)
				s.processTradeResult(result)
			})

			s.tradePairs.Store(tradeVariety.Symbol, engine)
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
	var data models_types.EventNotifyCancelOrder
	if err := json.Unmarshal(event.Message().Body, &data); err != nil {
		s.logger.Sugar().Errorf("matching notify cancel order unmarshal error: %v body: %s", err, string(event.Message().Body))
		return err
	}

	engine := s.engine(data.Symbol)
	if engine == nil {
		s.logger.Sugar().Errorf("matching engine not found for symbol: %s", data.Symbol)
		return nil
	}
	engine.RemoveItem(data.OrderSide, data.OrderId, matching_types.RemoveTypeByUser)
	return nil
}

func (s *Matching) engine(symbol string) *matching.Engine {
	if engine, ok := s.tradePairs.Load(symbol); ok {
		return engine.(*matching.Engine)
	}
	return nil
}

func (s *Matching) processCancelOrderResult(result matching_types.RemoveResult) {
	body, err := json.Marshal(result)
	if err != nil {
		s.logger.Sugar().Errorf("matching process cancel order result marshal error: %v", err)
		return
	}
	err = s.broker.Publish(context.Background(), models_types.TOPIC_PROCESS_ORDER_CANCEL, &broker.Message{
		Body: body,
	})
	if err != nil {
		s.logger.Sugar().Errorf("matching process cancel order result publish error: %v", err)
	}
}

func (s *Matching) processTradeResult(result matching_types.TradeResult) {
	body, err := result.MarshalBinary()
	if err != nil {
		s.logger.Sugar().Errorf("matching process trade result marshal error: %v", err)
		return
	}
	err = s.broker.Publish(context.Background(), models_types.TOPIC_ORDER_SETTLE, &broker.Message{
		Body: body,
	})
	if err != nil {
		s.logger.Sugar().Errorf("matching process trade result publish error: %v", err)
	}
}
