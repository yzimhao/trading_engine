package matching

import (
	"time"

	"github.com/shopspring/decimal"
)

func (e *Engine) processLimitOrder() {
	ok := func() bool {
		e.mx.Lock()
		defer e.mx.Unlock()

		if e.opts.pauseMatching {
			return false
		}
		if e.asks == nil || e.bids == nil {
			return false
		}

		if e.asks.Len() == 0 || e.bids.Len() == 0 {
			return false
		}

		askTop := e.asks.Top()
		bidTop := e.bids.Top()

		defer func() {
			if askTop.GetQuantity().Equal(decimal.Zero) {
				e.asks.Remove(askTop.GetUniqueId())
			}
			if bidTop.GetQuantity().Equal(decimal.Zero) {
				e.bids.Remove(bidTop.GetUniqueId())
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

			e.asks.SetQuantity(askTop, askTop.GetQuantity().Sub(curTradeQty))
			e.bids.SetQuantity(bidTop, bidTop.GetQuantity().Sub(curTradeQty))

			if askTop.GetCreateTime() >= bidTop.GetCreateTime() {
				curTradePrice = bidTop.GetPrice()
			} else {
				curTradePrice = askTop.GetPrice()
			}

			// 异步发送撮合结果，避免在持锁期阻塞
			e.emitTradeResult(e.tradeResult(askTop, bidTop, curTradePrice, curTradeQty, time.Now().UnixNano(), nil))
			return true
		} else {
			return false
		}

	}()

	if !ok {
		time.Sleep(time.Duration(100) * time.Millisecond)
	} else {
		if e.opts.debug {
			time.Sleep(time.Second * time.Duration(1))
		}
	}
}
