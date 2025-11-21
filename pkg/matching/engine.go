package matching

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/v2/pkg/matching/types"
	"go.uber.org/zap"
)

type Option func(opts *options)

type options struct {
	debug            bool
	priceDecimals    int32
	quantityDecimals int32
	pauseAcceptItem  bool
	pauseMatching    bool
	minTradeQuantity decimal.Decimal
	orderBookMaxLen  int
	logger           *zap.Logger
}

func defaultOptions() *options {
	return &options{
		debug:            false,
		priceDecimals:    2,
		quantityDecimals: 4,
		pauseAcceptItem:  false,
		pauseMatching:    false,
		orderBookMaxLen:  50,
		logger:           zap.NewNop(),
	}
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithPriceDecimals(decimals int32) Option {
	return func(opts *options) {
		opts.priceDecimals = decimals
	}
}

func WithQuantityDecimals(decimals int32) Option {
	return func(opts *options) {
		opts.quantityDecimals = decimals
	}
}

func WithDebug(debug bool) Option {
	return func(opts *options) {
		opts.debug = debug
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

type Engine struct {
	ctx          context.Context
	symbol       string
	opts         *options
	mx           sync.Mutex
	asks         *OrderQueue //orderqueue 不开放？？
	bids         *OrderQueue
	resultNotify chan types.TradeResult
	removeNotify chan types.RemoveResult
	logger       *zap.Logger

	onTradeResult  func(result types.TradeResult)
	onRemoveResult func(result types.RemoveResult)
}

func NewEngine(ctx context.Context, symbol string, opts ...Option) *Engine {
	o := defaultOptions()
	o.apply(opts...)

	e := &Engine{
		ctx:    ctx,
		symbol: symbol,
		opts:   o,
		asks:   NewQueue(),
		bids:   NewQueue(),
		// 提高缓冲，避免在高并发下阻塞撮合主循环
		resultNotify: make(chan types.TradeResult, 1024),
		removeNotify: make(chan types.RemoveResult, 1024),
		logger:       o.logger,
	}
	go e.matching()

	time.Sleep(time.Millisecond * 10) //???
	go e.orderBookTicker(e.asks)
	go e.orderBookTicker(e.bids)
	return e
}

func (e *Engine) Symbol() string {
	return e.symbol
}

func (e *Engine) PriceDecimals() int32 {
	return e.opts.priceDecimals
}

func (e *Engine) QuantityDecimals() int32 {
	return e.opts.quantityDecimals
}

func (e *Engine) SetPauseAcceptItem(pause bool) {
	e.opts.pauseAcceptItem = pause
}

func (e *Engine) SetPauseMatching(pause bool) {
	e.opts.pauseMatching = pause
}

func (e *Engine) AddItem(item QueueItem) error {
	e.mx.Lock()
	defer e.mx.Unlock()

	e.logger.Sugar().Debugf("[matching] AddItem %s", item.Marshal())

	if e.opts.pauseAcceptItem {
		return errors.New("engine is paused")
	}

	// 限价单：先主动撮合，剩余未成交部分才入订单簿
	if item.GetOrderType() == types.OrderTypeLimit {
		remain := item.GetQuantity()
		for remain.Cmp(decimal.Zero) > 0 {
			matched := false
			if item.GetOrderSide() == types.OrderSideSell {
				// 卖单主动撮合买队列
				if e.bids.Len() > 0 {
					bid := e.bids.Top()
					// 价格可成交
					if bid.GetPrice().Cmp(item.GetPrice()) >= 0 {
						tradeQty := decimal.Min(remain, bid.GetQuantity())
						// 成交价以买方（被动方）价格为准
						e.bids.SetQuantity(bid, bid.GetQuantity().Sub(tradeQty))
						remain = remain.Sub(tradeQty)
						e.emitTradeResult(e.tradeResult(item, bid, bid.GetPrice(), tradeQty, time.Now().UnixNano(), nil))
						if bid.GetQuantity().Equal(decimal.Zero) {
							e.bids.Remove(bid.GetUniqueId())
						}
						matched = true
					}
				}
			} else {
				// 买单主动撮合卖队列
				if e.asks.Len() > 0 {
					ask := e.asks.Top()
					// 价格可成交
					if item.GetPrice().Cmp(ask.GetPrice()) >= 0 {
						tradeQty := decimal.Min(remain, ask.GetQuantity())
						// 成交价以卖方（被动方）价格为准
						e.asks.SetQuantity(ask, ask.GetQuantity().Sub(tradeQty))
						remain = remain.Sub(tradeQty)
						e.emitTradeResult(e.tradeResult(ask, item, ask.GetPrice(), tradeQty, time.Now().UnixNano(), nil))
						if ask.GetQuantity().Equal(decimal.Zero) {
							e.asks.Remove(ask.GetUniqueId())
						}
						matched = true
					}
				}
			}
			if !matched {
				break // 无法继续撮合
			}
		}
		// 剩余未成交部分才入订单簿
		if remain.Cmp(decimal.Zero) > 0 {
			item.SetQuantity(remain)
			if item.GetOrderSide() == types.OrderSideSell {
				e.asks.Push(item)
			} else {
				e.bids.Push(item)
			}
		}
	} else {
		// 市价单逻辑保持原样
		if item.GetOrderSide() == types.OrderSideSell {
			go e.processMarketSell(item)
		} else {
			go e.processMarketBuy(item)
		}
	}

	return nil
}
func (e *Engine) RemoveItem(side types.OrderSide, unique string, removeType types.RemoveItemType) {
	e.mx.Lock()
	defer e.mx.Unlock()

	if side == types.OrderSideSell {
		e.asks.Remove(unique)
	} else if side == types.OrderSideBuy {
		e.bids.Remove(unique)
	} else {
		e.logger.Sugar().Warnf("removeItem %s side: %s unknown", unique, side)
	}

	e.emitRemoveResult(types.RemoveResult{
		Symbol:   e.symbol,
		UniqueId: unique,
		Type:     removeType,
	})
}

func (e *Engine) OnTradeResult(fn func(result types.TradeResult)) {
	e.onTradeResult = fn
}

func (e *Engine) OnRemoveResult(fn func(result types.RemoveResult)) {
	e.onRemoveResult = fn
}

func (e *Engine) GetAskOrderBook(size int) [][2]string {
	return e.orderBook(e.asks, size)
}

func (e *Engine) GetBidOrderBook(size int) [][2]string {
	return e.orderBook(e.bids, size)
}

func (e *Engine) AskQueue() *OrderQueue {
	return e.asks
}

func (e *Engine) BidQueue() *OrderQueue {
	return e.bids
}

func (e *Engine) Clean() {
	if e.opts.debug {
		e.mx.Lock()
		defer e.mx.Unlock()

		e.asks.clean()
		e.bids.clean()
	}
}

func (e *Engine) matching() {
	e.logger.Debug("[matching] start...")
	for {
		select {
		case <-e.ctx.Done():
			return
		case result := <-e.resultNotify:
			e.logger.Debug("[matching] trade result", zap.Any("result", result))
			if e.onTradeResult != nil {
				e.onTradeResult(result)
			}
		case result := <-e.removeNotify:
			e.logger.Debug("[matching] cancel item", zap.Any("result", result))
			if e.onRemoveResult != nil {
				e.onRemoveResult(result)
			}
		default:
			// e.logger.Debug("[matching] processLimitOrder")
			e.processLimitOrder()
		}
	}
}

// emitTradeResult 异步发送撮合结果，避免在持锁区或处理路径上阻塞
func (e *Engine) emitTradeResult(tr types.TradeResult) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				e.logger.Sugar().Errorf("emitTradeResult panic: %v", r)
			}
		}()
		e.resultNotify <- tr
	}()
}

// emitRemoveResult 异步发送移除/撤单通知
func (e *Engine) emitRemoveResult(rr types.RemoveResult) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				e.logger.Sugar().Errorf("emitRemoveResult panic: %v", r)
			}
		}()
		e.removeNotify <- rr
	}()
}

func (e *Engine) tradeResult(ask, bid QueueItem, price, tradeQty decimal.Decimal, tradeAt int64, marketOrder *types.MarketOrderInfo) types.TradeResult {

	tradeResult := types.TradeResult{
		Symbol:          e.symbol,
		AskOrderId:      ask.GetUniqueId(),
		BidOrderId:      bid.GetUniqueId(),
		TradeQuantity:   tradeQty,
		TradePrice:      price,
		TradeTime:       tradeAt,
		MarketOrderInfo: marketOrder,
	}

	if ask.GetCreateTime() < bid.GetCreateTime() {
		tradeResult.TradeBy = types.TradeByBuyer
	} else {
		tradeResult.TradeBy = types.TradeBySeller
	}

	// if tradeAt > e.latestPriceAt {
	// 	t.latestPrice = price
	// 	t.latestPriceAt = trade_at
	// }

	return tradeResult
}
