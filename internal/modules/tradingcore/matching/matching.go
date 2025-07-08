package matching

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/duolacloud/crud-core/cache"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	notification_ws "github.com/yzimhao/trading_engine/v2/internal/modules/notification/ws"
	"github.com/yzimhao/trading_engine/v2/internal/modules/tradingcore/orderlock"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	"github.com/yzimhao/trading_engine/v2/internal/types"
	models_types "github.com/yzimhao/trading_engine/v2/internal/types"
	"github.com/yzimhao/trading_engine/v2/pkg/matching"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	CacheKeyOrderbook = "orderbook.%s" //example: orderbook.btcusdt
)

type inContext struct {
	fx.In
	Produce     *provider.Produce
	Consume     *provider.Consume
	Logger      *zap.Logger
	ProductRepo persistence.ProductRepository
	Viper       *viper.Viper
	Cache       cache.Cache
	OrderRepo   persistence.OrderRepository
	Ws          *notification_ws.WsManager
	Locker      *orderlock.OrderLock
}

type Matching struct {
	produce     *provider.Produce
	consume     *provider.Consume
	logger      *zap.Logger
	productRepo persistence.ProductRepository
	tradePairs  sync.Map
	viper       *viper.Viper
	cache       cache.Cache
	orderRepo   persistence.OrderRepository
	ws          *notification_ws.WsManager
	locker      *orderlock.OrderLock
}

func NewMatching(in inContext) *Matching {
	return &Matching{
		produce:     in.Produce,
		consume:     in.Consume,
		logger:      in.Logger,
		productRepo: in.ProductRepo,
		viper:       in.Viper,
		cache:       in.Cache,
		orderRepo:   in.OrderRepo,
		ws:          in.Ws,
		locker:      in.Locker,
	}
}

func (s *Matching) InitEngine() {
	s.logger.Sugar().Infof("init matching engine")
	localSymbols := s.viper.GetStringSlice("matching.local_symbols")

	// load trade pair

	var (
		products []entities.Product
	)

	if err := s.productRepo.DB().Model(entities.Product{}).Find(&products).Error; err != nil {
		s.logger.Sugar().Errorf("query trade product error: %v", err)
		return
	}

	for _, product := range products {
		if len(localSymbols) > 0 {
			if !slices.Contains(localSymbols, product.Symbol) {
				continue
			}
		}

		opts := []matching.Option{
			matching.WithPriceDecimals(int32(product.PriceDecimals)),
			matching.WithQuantityDecimals(int32(product.QtyDecimals)),
			matching.WithLogger(s.logger),
		}
		engine := matching.NewEngine(context.Background(), product.Symbol, opts...)

		engine.OnRemoveResult(func(result matching_types.RemoveResult) {
			s.logger.Sugar().Infof("symbol: %s remove result: %+v", result.Symbol, result)
			// time.Sleep(time.Second) //这里如何延迟取消订单的通知？
			s.processCancelOrderResult(result)
		})
		engine.OnTradeResult(func(result matching_types.TradeResult) {
			s.logger.Sugar().Infof("symbol: %s trade result: %+v", result.Symbol, result)
			s.processTradeResult(result)
		})

		s.tradePairs.Store(product.Symbol, engine)
		s.logger.Sugar().Infof("init matching engine for symbol: %s", product.Symbol)

		go s.flushOrderbookToCache(context.Background(), product.Symbol)

		//TODO  load order from db
		s.loadUnfinishedOrders(context.Background(), product.Symbol)
	}

}

func (s *Matching) Subscribe() {
	s.consume.Subscribe(models_types.TOPIC_ORDER_NEW, func(ctx context.Context, data []byte) {
		s.OnNewOrder(ctx, data)
	})
	s.consume.Subscribe(models_types.TOPIC_NOTIFY_ORDER_CANCEL, func(ctx context.Context, data []byte) {
		s.OnNotifyCancelOrder(ctx, data)
	})
}

func (s *Matching) OnNewOrder(ctx context.Context, msg []byte) error {
	s.logger.Sugar().Debugf("listen new order %s", msg)

	var order models_types.EventOrderNew
	if err := json.Unmarshal(msg, &order); err != nil {
		s.logger.Sugar().Errorf("matching new order unmarshal error: %v body: %s", err, msg)
		return err
	}

	s.logger.Sugar().Debugf("listen new order: %+v", order)

	var item matching.QueueItem
	if order.OrderType == matching_types.OrderTypeLimit {
		if order.OrderSide == matching_types.OrderSideSell {
			item = matching.NewAskLimitItem(order.OrderId, order.Price, order.Quantity, order.NanoTime)
		} else {
			item = matching.NewBidLimitItem(order.OrderId, order.Price, order.Quantity, order.NanoTime)
		}
	} else if order.OrderType == matching_types.OrderTypeMarket {
		// 按成交金额
		if order.Amount.Cmp(decimal.Zero) > 0 {
			if order.OrderSide == matching_types.OrderSideSell {
				item = matching.NewAskMarketAmountItem(order.OrderId, order.Amount, order.MaxQty, order.NanoTime)
			} else {
				item = matching.NewBidMarketAmountItem(order.OrderId, order.Amount, order.NanoTime)
			}
		} else {
			// 按成交量
			if order.OrderSide == matching_types.OrderSideSell {
				item = matching.NewAskMarketQtyItem(order.OrderId, order.MaxQty, order.NanoTime)
			} else {
				item = matching.NewBidMarketQtyItem(order.OrderId, order.Quantity, order.MaxAmount, order.NanoTime)
			}
		}
	}

	if engine := s.engine(order.Symbol); engine != nil {
		engine.AddItem(item)
		s.logger.Sugar().Debugf("add item to engine %s, askLen: %d, bidLen: %d", order.Symbol, engine.AskQueue().Len(), engine.BidQueue().Len())
	}
	return nil
}

func (s *Matching) OnNotifyCancelOrder(ctx context.Context, msg []byte) error {
	var data models_types.EventNotifyCancelOrder
	if err := json.Unmarshal(msg, &data); err != nil {
		s.logger.Sugar().Errorf("matching notify cancel order unmarshal error: %v body: %s", err, msg)
		return err
	}

	engine := s.engine(data.Symbol)
	if engine == nil {
		s.logger.Sugar().Errorf("matching engine not found for symbol: %s", data.Symbol)
		return nil
	}
	engine.RemoveItem(data.OrderSide, data.OrderId, data.Type)
	return nil
}

func (s *Matching) engine(symbol string) *matching.Engine {
	if engine, ok := s.tradePairs.Load(symbol); ok {
		return engine.(*matching.Engine)
	}
	return nil
}

func (s *Matching) processCancelOrderResult(result matching_types.RemoveResult) {
	data := types.EventCancelOrder{
		Symbol:  result.Symbol,
		OrderId: result.UniqueId,
	}
	body, err := json.Marshal(data)
	if err != nil {
		s.logger.Sugar().Errorf("matching process cancel order result marshal error: %v", err)
		return
	}

	err = s.produce.Publish(context.Background(), models_types.TOPIC_PROCESS_ORDER_CANCEL, body)
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
	//处理成交结果之前，对相关订单加锁
	ctx := context.Background()
	if err := s.locker.Lock(ctx, result.AskOrderId, result.BidOrderId); err != nil {
		s.logger.Sugar().Errorf("matching process trade result lock error: %v", err)
		return
	}

	err = s.produce.Publish(ctx, models_types.TOPIC_ORDER_SETTLE, body)
	if err != nil {
		s.logger.Sugar().Errorf("matching process trade result publish error: %v", err)
	}
}

func (s *Matching) flushOrderbookToCache(ctx context.Context, symbol string) {
	ticker := time.NewTicker(100 * time.Millisecond)

	engine := s.engine(symbol)
	if engine == nil {
		s.logger.Sugar().Errorf("matching engine not found for symbol: %s", symbol)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			asks := engine.GetAskOrderBook(10)
			bids := engine.GetBidOrderBook(10)

			data := map[string]any{
				"asks": asks,
				"bids": bids,
			}
			err := s.cache.Set(ctx, fmt.Sprintf(CacheKeyOrderbook, engine.Symbol()), data, cache.WithExpiration(time.Second*5))
			if err != nil {
				s.logger.Sugar().Errorf("matching flush orderbook to cache error: %v", err)
			}

			//broadcast depth data
			if err := s.ws.Broadcast(ctx, notification_ws.MsgDepthTpl.Format(map[string]string{"symbol": engine.Symbol()}), data); err != nil {
				s.logger.Sugar().Errorf("matching ws broadcast error: %v", err)
			}
		}
	}
}

func (s *Matching) loadUnfinishedOrders(ctx context.Context, symbol string) error {
	orders, err := s.orderRepo.LoadUnfinishedOrders(ctx, symbol)
	if err != nil {
		s.logger.Sugar().Errorf("matching load unfinished orders error: %v", err)
		return err
	}

	for _, order := range orders {
		var item matching.QueueItem
		if order.OrderType == matching_types.OrderTypeLimit {
			if order.OrderSide == matching_types.OrderSideSell {
				item = matching.NewAskLimitItem(order.OrderId, order.Price, order.Quantity.Sub(order.FinishedQty), order.NanoTime)
			} else {
				item = matching.NewBidLimitItem(order.OrderId, order.Price, order.Quantity.Sub(order.FinishedQty), order.NanoTime)
			}
			if engine := s.engine(order.Symbol); engine != nil {

				engine.AddItem(item)
				s.logger.Sugar().Debugf("load unfinished order %s to engine %s, askLen: %d, bidLen: %d", order.OrderId, order.Symbol, engine.AskQueue().Len(), engine.BidQueue().Len())
			}
		}

	}

	return nil
}
