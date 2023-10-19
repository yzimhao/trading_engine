package trading_core

import (
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

type TradePair struct {
	Symbol         string
	ChTradeResult  chan TradeResult
	ChCancelResult chan string

	pausePushNew bool
	pauseMatch   bool
	triggerEvent bool

	onEventTrade func(TradeResult)

	priceDigit    int
	quantityDigit int
	miniTradeQty  decimal.Decimal

	latestPriceAt int64
	latestPrice   decimal.Decimal

	askQueue *OrderQueue
	bidQueue *OrderQueue

	w sync.Mutex
}

func NewTradePair(symbol string, priceDigit, quantityDigit int) *TradePair {
	t := &TradePair{
		Symbol:         symbol,
		ChTradeResult:  make(chan TradeResult, 10),
		ChCancelResult: make(chan string, 10),

		pausePushNew: false,
		pauseMatch:   false,

		priceDigit:    priceDigit,
		quantityDigit: quantityDigit,
		miniTradeQty:  decimal.New(1, int32(-quantityDigit)),

		askQueue: NewQueue(),
		bidQueue: NewQueue(),
	}

	go t.depthTicker(t.askQueue)
	go t.depthTicker(t.bidQueue)
	go t.matching()
	return t
}

func (t *TradePair) OnEvent(fn_update, fn_remove func(QueueItem), fn_trade func(TradeResult)) {
	t.askQueue.onEventRemove = fn_remove
	t.askQueue.onEventUpdate = fn_update

	t.bidQueue.onEventRemove = fn_remove
	t.bidQueue.onEventUpdate = fn_update

	t.onEventTrade = fn_trade
}

func (t *TradePair) IsPausePushNew() bool {
	return t.pausePushNew
}

func (t *TradePair) SetPausePushNew(v bool) {
	t.pausePushNew = v
}

func (t *TradePair) SetPauseMatch(v bool) {
	t.pauseMatch = v
}

func (t *TradePair) SetTriggerEvent(v bool) {
	t.triggerEvent = v
}

func (t *TradePair) TriggerEvent() bool {
	return t.triggerEvent
}

func (t *TradePair) PushNewOrder(item QueueItem) {
	t.handlerNewOrder(item)
}

func (t *TradePair) CancelOrder(side OrderSide, uniq string) {
	if side == OrderSideSell {
		t.askQueue.Remove(uniq)
	} else {
		t.bidQueue.Remove(uniq)
	}
	//删除成功后需要发送通知
	t.ChCancelResult <- uniq
}

func (t *TradePair) AskLen() int {
	t.w.Lock()
	defer t.w.Unlock()

	return t.askQueue.Len()
}

func (t *TradePair) BidLen() int {
	t.w.Lock()
	defer t.w.Unlock()

	return t.bidQueue.Len()
}

func (t *TradePair) LatestPrice() (decimal.Decimal, int64) {
	t.w.Lock()
	defer t.w.Unlock()
	return t.latestPrice, t.latestPriceAt
}

func (t *TradePair) SetLatestPrice(p decimal.Decimal) {
	t.w.Lock()
	defer t.w.Unlock()
	t.latestPrice = p
}

func (t *TradePair) cleanAll() {
	t.w.Lock()
	defer t.w.Unlock()

	//同时清空两个队列
	t.askQueue.clean()
	t.bidQueue.clean()
}

func (t *TradePair) matching() {

	for {
		select {
		default:
			t.handlerLimitOrder()

		}

	}

}

func (t *TradePair) handlerNewOrder(newOrder QueueItem) {
	t.w.Lock()
	defer t.w.Unlock()

	if t.pausePushNew {
		return
	}

	if newOrder.GetOrderType() == OrderTypeLimit {
		if newOrder.GetOrderSide() == OrderSideSell {
			t.askQueue.Push(newOrder)
		} else {
			t.bidQueue.Push(newOrder)
		}
	} else {
		//市价单处理
		if newOrder.GetOrderSide() == OrderSideSell {
			t.doMarketSell(newOrder)
		} else {
			t.doMarketBuy(newOrder)
		}
	}

}

func (t *TradePair) handlerLimitOrder() {
	ok := func() bool {
		t.w.Lock()
		defer t.w.Unlock()

		if t.pauseMatch {
			return false
		}
		if t.askQueue == nil || t.bidQueue == nil {
			return false
		}

		if t.askQueue.Len() == 0 || t.bidQueue.Len() == 0 {
			return false
		}

		askTop := t.askQueue.Top()
		bidTop := t.bidQueue.Top()

		defer func() {
			if askTop.GetQuantity().Equal(decimal.Zero) {
				t.askQueue.Remove(askTop.GetUniqueId())
			}
			if bidTop.GetQuantity().Equal(decimal.Zero) {
				t.bidQueue.Remove(bidTop.GetUniqueId())
			}
		}()

		if bidTop.GetPrice().Cmp(askTop.GetPrice()) >= 0 {
			curTradeQty := decimal.Zero
			var curTradePrice decimal.Decimal
			if bidTop.GetQuantity().Cmp(askTop.GetQuantity()) >= 0 {
				curTradeQty = askTop.GetQuantity()
			} else if bidTop.GetQuantity().Cmp(askTop.GetQuantity()) == -1 {
				curTradeQty = bidTop.GetQuantity()
			}

			t.askQueue.SetQuantity(askTop, askTop.GetQuantity().Sub(curTradeQty))
			t.bidQueue.SetQuantity(bidTop, bidTop.GetQuantity().Sub(curTradeQty))

			if askTop.GetCreateTime() >= bidTop.GetCreateTime() {
				curTradePrice = bidTop.GetPrice()
			} else {
				curTradePrice = askTop.GetPrice()
			}

			go t.sendTradeResultNotify(askTop, bidTop, curTradePrice, curTradeQty, time.Now().UnixNano(), "")
			return true
		} else {
			return false
		}

	}()

	if !ok {
		time.Sleep(time.Duration(60) * time.Millisecond)
	} else {
		if Debug {
			time.Sleep(time.Second * time.Duration(1))
		}
	}
}

func (t *TradePair) doMarketBuy(item QueueItem) {

	for {
		ok := func() bool {

			if t.askQueue.Len() == 0 {
				return false
			}

			ask := t.askQueue.Top()
			if item.GetOrderType() == OrderTypeMarketQuantity {
				maxQty := func(remainAmount, marketPrice, needQty decimal.Decimal) decimal.Decimal {
					qty := remainAmount.Div(marketPrice)
					return decimal.Min(qty, needQty)
				}
				maxTradeQty := maxQty(item.GetAmount(), ask.GetPrice(), item.GetQuantity())
				curTradeQty := decimal.Zero

				//市价按买入数量
				if maxTradeQty.Cmp(t.miniTradeQty) < 0 {
					return false
				}

				if ask.GetQuantity().Cmp(maxTradeQty) <= 0 {
					curTradeQty = ask.GetQuantity()
					t.askQueue.Remove(ask.GetUniqueId())
				} else {
					curTradeQty = maxTradeQty
					t.askQueue.SetQuantity(ask, ask.GetQuantity().Sub(curTradeQty))
				}

				item.SetQuantity(item.GetQuantity().Sub(curTradeQty))
				item.SetAmount(item.GetAmount().Sub(curTradeQty.Mul(ask.GetPrice())))

				//检查本次循环撮合是否是该订单最后一次撮合
				//如果是则标记该市价订单已经完成了
				//结束的条件：
				// a.对面订单列表空了
				// b.已经达到了用户需要的数量
				// c.剩余资金已经不能达到最小成交需求
				if t.askQueue.Len() == 0 || item.GetQuantity().Equal(decimal.Zero) ||
					maxQty(item.GetAmount(), t.askQueue.Top().GetPrice(), item.GetQuantity()).Cmp(t.miniTradeQty) < 0 {
					go t.sendTradeResultNotify(ask, item, ask.GetPrice(), curTradeQty, time.Now().UnixNano(), item.GetUniqueId())
				} else {
					go t.sendTradeResultNotify(ask, item, ask.GetPrice(), curTradeQty, time.Now().UnixNano(), "")
				}

				return true
			} else if item.GetOrderType() == OrderTypeMarketAmount {
				//市价-按成交金额
				//成交金额不包含手续费，手续费应该由上层系统计算提前预留
				//撮合会针对这个金额最大限度的买入
				if ask.GetPrice().Cmp(decimal.Zero) <= 0 {
					return false
				}

				maxQty := func(amount, price decimal.Decimal) decimal.Decimal {
					return amount.Div(price)
				}

				maxTradeQty := maxQty(item.GetAmount(), ask.GetPrice()) //item.GetAmount().Div(ask.GetPrice())
				curTradeQty := decimal.Zero

				if maxTradeQty.Cmp(t.miniTradeQty) < 0 {
					return false
				}
				if ask.GetQuantity().Cmp(maxTradeQty) <= 0 {
					curTradeQty = ask.GetQuantity()
					t.askQueue.Remove(ask.GetUniqueId())
				} else {
					curTradeQty = maxTradeQty
					// ask.SetQuantity(ask.GetQuantity().Sub(curTradeQty))
					t.askQueue.SetQuantity(ask, ask.GetQuantity().Sub(curTradeQty))
				}

				//部分成交了，需要更新这个单的剩余可成交金额，用于下一轮重新计算最大成交量
				item.SetAmount(item.GetAmount().Sub(curTradeQty.Mul(ask.GetPrice())))
				item.SetQuantity(item.GetQuantity().Add(curTradeQty))

				//检查本次循环撮合是否是该订单最后一次撮合
				//结束的条件：
				// a.对面订单列表空了
				// b.已经达到了用户需要的数量
				// c.剩余资金已经不能达到最小成交需求
				if t.askQueue.Len() == 0 || item.GetQuantity().Equal(decimal.Zero) ||
					maxQty(item.GetAmount(), t.askQueue.Top().GetPrice()).Cmp(t.miniTradeQty) < 0 {
					go t.sendTradeResultNotify(ask, item, ask.GetPrice(), curTradeQty, time.Now().UnixNano(), item.GetUniqueId())
				} else {
					go t.sendTradeResultNotify(ask, item, ask.GetPrice(), curTradeQty, time.Now().UnixNano(), "")
				}
				return true
			}

			return false
		}()

		if !ok {
			break
		}

	}
}
func (t *TradePair) doMarketSell(item QueueItem) {

	for {
		ok := func() bool {

			if t.bidQueue.Len() == 0 {
				return false
			}

			bid := t.bidQueue.Top()
			if item.GetOrderType() == OrderTypeMarketQuantity {

				curTradeQuantity := decimal.Zero
				//市价按买入数量
				if item.GetQuantity().Equal(decimal.Zero) {
					return false
				}

				if bid.GetQuantity().Cmp(item.GetQuantity()) <= 0 {
					curTradeQuantity = bid.GetQuantity()
					t.bidQueue.Remove(bid.GetUniqueId())
				} else {
					curTradeQuantity = item.GetQuantity()
					// bid.SetQuantity(bid.GetQuantity().Sub(curTradeQuantity))
					t.bidQueue.SetQuantity(bid, bid.GetQuantity().Sub(curTradeQuantity))
				}

				item.SetQuantity(item.GetQuantity().Sub(curTradeQuantity))

				//退出条件
				// a.对面订单空了
				// b.市价订单完全成交了
				if t.bidQueue.Len() == 0 || item.GetQuantity().Equal(decimal.Zero) {
					go t.sendTradeResultNotify(item, bid, bid.GetPrice(), curTradeQuantity, time.Now().UnixNano(), item.GetUniqueId())
				} else {
					go t.sendTradeResultNotify(item, bid, bid.GetPrice(), curTradeQuantity, time.Now().UnixNano(), "")
				}

				return true
			} else if item.GetOrderType() == OrderTypeMarketAmount {
				//市价-按成交金额成交
				if bid.GetPrice().Cmp(decimal.Zero) <= 0 {
					return false
				}

				maxQty := func(amount, price, needQty decimal.Decimal) decimal.Decimal {
					a := amount.Div(price).Truncate(int32(t.quantityDigit))
					return decimal.Min(a, needQty)
				}
				maxTradeQty := maxQty(item.GetAmount(), bid.GetPrice(), item.GetQuantity())

				curTradeQty := decimal.Zero
				if maxTradeQty.Cmp(t.miniTradeQty) < 0 {
					return false
				}

				if bid.GetQuantity().Cmp(maxTradeQty) <= 0 {
					curTradeQty = bid.GetQuantity()
					t.bidQueue.Remove(bid.GetUniqueId())
				} else {
					curTradeQty = maxTradeQty
					// bid.SetQuantity(bid.GetQuantity().Sub(curTradeQty))
					t.bidQueue.SetQuantity(bid, bid.GetQuantity().Sub(curTradeQty))
				}

				item.SetAmount(item.GetAmount().Sub(curTradeQty.Mul(bid.GetPrice())))
				//市价 按成交额卖出时，需要用户持有的资产数量来进行限制
				item.SetQuantity(item.GetQuantity().Sub(curTradeQty))

				//退出条件
				// a.对面订单空了
				// b.金额完全成交
				// c.剩余资金不满足最小成交量
				if t.bidQueue.Len() == 0 || maxQty(item.GetAmount(), t.bidQueue.Top().GetPrice(), item.GetQuantity()).Cmp(t.miniTradeQty) < 0 {
					go t.sendTradeResultNotify(item, bid, bid.GetPrice(), curTradeQty, time.Now().UnixNano(), item.GetUniqueId())
				} else {
					go t.sendTradeResultNotify(item, bid, bid.GetPrice(), curTradeQty, time.Now().UnixNano(), "")
				}

				return true
			}

			return false
		}()

		if !ok {
			break
		}

	}
}

func (t *TradePair) sendTradeResultNotify(ask, bid QueueItem, price, tradeQty decimal.Decimal, trade_at int64, last string) {
	// t.w.Lock()
	// defer t.w.Unlock()

	tradelog := TradeResult{}
	tradelog.Symbol = t.Symbol
	tradelog.AskOrderId = ask.GetUniqueId()
	tradelog.BidOrderId = bid.GetUniqueId()
	tradelog.TradeQuantity = tradeQty
	tradelog.TradePrice = price
	tradelog.TradeTime = trade_at //精确到纳秒
	tradelog.Last = last          //标记市价订单已经完成

	if trade_at > t.latestPriceAt {
		t.latestPrice = price
		t.latestPriceAt = trade_at
	}

	if Debug {
		// logrus.Infof("%s tradelog: %+v", t.Symbol, tradelog)
	}

	t.ChTradeResult <- tradelog

	if t.onEventTrade != nil {
		t.onEventTrade(tradelog)
	}
}
