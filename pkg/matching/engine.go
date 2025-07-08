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
		ctx:          ctx,
		symbol:       symbol,
		opts:         o,
		asks:         NewQueue(),
		bids:         NewQueue(),
		resultNotify: make(chan types.TradeResult, 1),
		removeNotify: make(chan types.RemoveResult, 1),
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

	if item.GetOrderType() == types.OrderTypeLimit {
		if item.GetOrderSide() == types.OrderSideSell {
			e.asks.Push(item)
		} else {
			e.bids.Push(item)
		}
	} else {
		if item.GetOrderSide() == types.OrderSideSell {
			go e.processMarketSell(item)
		} else {
			go e.processMarketBuy(item)
		}
	}

	return nil
}
func (e *Engine) RemoveItem(side types.OrderSide, unique string, removeType types.RemoveType) {
	e.mx.Lock()
	defer e.mx.Unlock()

	if side == types.OrderSideSell {
		e.asks.Remove(unique)
	} else {
		e.bids.Remove(unique)
	}

	e.removeNotify <- types.RemoveResult{
		Symbol:   e.symbol,
		UniqueId: unique,
		Type:     removeType,
	}
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

func (e *Engine) tradeResult(ask, bid QueueItem, price, tradeQty decimal.Decimal, tradeAt int64, remainder string) types.TradeResult {

	tradeResult := types.TradeResult{
		Symbol:                 e.symbol,
		AskOrderId:             ask.GetUniqueId(),
		BidOrderId:             bid.GetUniqueId(),
		TradeQuantity:          tradeQty,
		TradePrice:             price,
		TradeTime:              tradeAt,
		RemainderMarketOrderId: remainder,
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
