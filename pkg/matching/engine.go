package matching

import (
	"context"
	"errors"
	"sync"

	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/v2/pkg/matching/types"
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
}

func defaultOptions() *options {
	return &options{
		debug:            false,
		priceDecimals:    2,
		quantityDecimals: 4,
		pauseAcceptItem:  false,
		pauseMatching:    false,
		orderBookMaxLen:  50,
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

type Engine struct {
	ctx          context.Context
	symbol       string
	opts         *options
	mx           sync.Mutex
	asks         *OrderQueue //orderqueue 不开放？？
	bids         *OrderQueue
	resultNotify chan types.TradeResult
	removeNotify chan types.RemoveResult

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
		resultNotify: make(chan types.TradeResult, 10),
		removeNotify: make(chan types.RemoveResult, 10),
	}
	go e.matching()
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
			e.processMarketSell(item)
		} else {
			e.processMarketBuy(item)
		}
	}

	return nil
}
func (e *Engine) RemoveItem(side types.OrderSide, unique string) {
	e.mx.Lock()
	defer e.mx.Unlock()

	if side == types.OrderSideSell {
		e.asks.Remove(unique)
	} else {
		e.bids.Remove(unique)
	}

	e.removeNotify <- types.RemoveResult{
		UniqueId: unique,
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
	for {
		select {
		case <-e.ctx.Done():
			return
		case result := <-e.resultNotify:
			if e.onTradeResult != nil {
				e.onTradeResult(result)
			}
		case result := <-e.removeNotify:
			if e.onRemoveResult != nil {
				e.onRemoveResult(result)
			}
		default:
			e.processLimitOrder()
		}
	}
}

func (e *Engine) tradeResult(ask, bid QueueItem, price, tradeQty decimal.Decimal, tradeAt int64, remainder string) types.TradeResult {

	tradeResult := types.TradeResult{
		Symbol:                 e.symbol,
		AskOrderId:             ask.GetUniqueId(),
		BidOrderId:             bid.GetUniqueId(),
		TradeQuantity:          tradeQty.String(),
		TradePrice:             price.String(),
		TradeTime:              tradeAt,
		RemainderMarketOrderId: remainder,
	}

	if ask.GetCreateTime() < bid.GetCreateTime() {
		tradeResult.TradeBy = types.ByBuyer
	} else {
		tradeResult.TradeBy = types.BySeller
	}

	// if tradeAt > e.latestPriceAt {
	// 	t.latestPrice = price
	// 	t.latestPriceAt = trade_at
	// }

	return tradeResult
}
