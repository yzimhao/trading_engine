package matching

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/v2/pkg/matching/types"
	"go.uber.org/zap"
)

func (e *Engine) processMarketBuy(item QueueItem) {
	e.mx.Lock()
	defer e.mx.Unlock()

	trade_cnt := 0
	for {
		ok := func() bool {

			if e.asks.Len() == 0 {
				return false
			}

			ask := e.asks.Top()

			if item.GetSubOrderType() == types.SubOrderTypeMarketByQty {
				maxQty := func(remainAmount, marketPrice, needQty decimal.Decimal) decimal.Decimal {
					qty := remainAmount.Div(marketPrice)
					return decimal.Min(qty, needQty).Truncate(e.opts.quantityDecimals)
				}

				e.logger.Sugar().Debugf("[matching] processMarketBuy item: %s", item.Marshal())

				maxTradeQty := maxQty(item.GetHoldAmount(), ask.GetPrice(), item.GetQuantity())

				curTradeQty := decimal.Zero

				//市价按买入数量
				if maxTradeQty.Cmp(e.opts.minTradeQuantity) < 0 {
					return false
				}

				e.logger.Debug("[matching] processMarketBuy", zap.String("maxTradeQty", maxTradeQty.String()),
					zap.String("ask.GetQuantity()", ask.GetQuantity().String()),
					zap.Int("trade_cnt", trade_cnt),
				)

				if ask.GetQuantity().Cmp(maxTradeQty) <= 0 {
					curTradeQty = ask.GetQuantity()
					e.asks.Remove(ask.GetUniqueId())
				} else {
					curTradeQty = maxTradeQty
					e.asks.SetQuantity(ask, ask.GetQuantity().Sub(curTradeQty))
				}

				if curTradeQty.Equal(decimal.Zero) {
					return false
				}

				item.SetQuantity(item.GetQuantity().Sub(curTradeQty))
				item.SetHoldAmount(item.GetHoldAmount().Sub(curTradeQty.Mul(ask.GetPrice())))

				//检查本次循环撮合是否是该订单最后一次撮合
				//如果是则标记该市价订单已经完成了
				//结束的条件：
				// a.对面订单列表空了
				// b.已经达到了用户需要的数量
				// c.剩余资金已经不能达到最小成交需求
				if e.asks.Len() == 0 || item.GetQuantity().Equal(decimal.Zero) ||
					maxQty(item.GetHoldAmount(), e.asks.Top().GetPrice(), item.GetQuantity()).Cmp(e.opts.minTradeQuantity) <= 0 {
					e.resultNotify <- e.tradeResult(ask, item, ask.GetPrice(), curTradeQty, time.Now().UnixNano(), &types.MarketOrderInfo{
						OrderId:      item.GetUniqueId(),
						IsFinalTrade: true,
					})
				} else {
					e.resultNotify <- e.tradeResult(ask, item, ask.GetPrice(), curTradeQty, time.Now().UnixNano(), &types.MarketOrderInfo{
						OrderId:      item.GetUniqueId(),
						IsFinalTrade: false,
					})
				}

				return true
			} else if item.GetSubOrderType() == types.SubOrderTypeMarketByAmount {
				//市价-按成交金额
				//成交金额不包含手续费，手续费应该由上层系统计算提前预留
				//撮合会针对这个金额最大限度的买入
				if ask.GetPrice().Cmp(decimal.Zero) <= 0 {
					return false
				}

				// 持有的金额/卖家报价=最大可以买入的数量
				maxQty := func(holdAmount, price decimal.Decimal) decimal.Decimal {
					return holdAmount.Div(price).Truncate(e.opts.quantityDecimals)
				}

				maxTradeQty := maxQty(item.GetHoldAmount(), ask.GetPrice())
				curTradeQty := decimal.Zero

				if maxTradeQty.Cmp(e.opts.minTradeQuantity) < 0 {
					return false
				}

				if ask.GetQuantity().Cmp(maxTradeQty) <= 0 {
					curTradeQty = ask.GetQuantity()
					e.asks.Remove(ask.GetUniqueId())
				} else {
					curTradeQty = maxTradeQty
					e.asks.SetQuantity(ask, ask.GetQuantity().Sub(curTradeQty))
				}

				if curTradeQty.Equal(decimal.Zero) {
					return false
				}

				//部分成交了，需要更新这个单的剩余可成交金额，用于下一轮重新计算最大成交量
				item.SetHoldAmount(item.GetHoldAmount().Sub(curTradeQty.Mul(ask.GetPrice())))
				item.SetQuantity(item.GetQuantity().Add(curTradeQty))

				//检查本次循环撮合是否是该订单最后一次撮合
				//结束的条件：
				// a.对面订单列表空了
				// b.已经达到了用户需要的数量
				// c.剩余资金已经不能达到最小成交需求
				if e.asks.Len() == 0 || item.GetQuantity().Equal(decimal.Zero) ||
					maxQty(item.GetHoldAmount(), e.asks.Top().GetPrice()).Cmp(e.opts.minTradeQuantity) <= 0 {
					e.resultNotify <- e.tradeResult(ask, item, ask.GetPrice(), curTradeQty, time.Now().UnixNano(), &types.MarketOrderInfo{
						OrderId:      item.GetUniqueId(),
						IsFinalTrade: true,
					})
				} else {
					e.resultNotify <- e.tradeResult(ask, item, ask.GetPrice(), curTradeQty, time.Now().UnixNano(), &types.MarketOrderInfo{
						OrderId:      item.GetUniqueId(),
						IsFinalTrade: false,
					})
				}
				return true
			}

			return false
		}()

		if !ok {
			//TODO 市价订单，需要触发一个自动取消订单操作 这里的取消可能会比成交记录先到到处理端
			go func() {
				// time.Sleep(time.Second)
				e.removeNotify <- types.RemoveResult{
					Symbol:   e.symbol,
					UniqueId: item.GetUniqueId(),
					Type:     types.RemoveItemTypeByMarket,
				}
			}()
			break
		} else {
			trade_cnt++
		}

	}
}
func (e *Engine) processMarketSell(item QueueItem) {
	e.mx.Lock()
	defer e.mx.Unlock()

	trade_cnt := 0
	for {
		ok := func() bool {

			if e.bids.Len() == 0 {
				return false
			}

			bid := e.bids.Top()
			if item.GetSubOrderType() == types.SubOrderTypeMarketByQty {

				curTradeQuantity := decimal.Zero
				//市价按买入数量
				if item.GetQuantity().Equal(decimal.Zero) {
					return false
				}

				if bid.GetQuantity().Cmp(item.GetQuantity()) <= 0 {
					curTradeQuantity = bid.GetQuantity()
					e.bids.Remove(bid.GetUniqueId())
				} else {
					curTradeQuantity = item.GetQuantity()
					e.bids.SetQuantity(bid, bid.GetQuantity().Sub(curTradeQuantity))
				}

				item.SetQuantity(item.GetQuantity().Sub(curTradeQuantity))

				//退出条件
				// a.对面订单空了
				// b.市价订单完全成交了
				if e.bids.Len() == 0 || item.GetQuantity().Equal(decimal.Zero) {
					e.resultNotify <- e.tradeResult(item, bid, bid.GetPrice(), curTradeQuantity, time.Now().UnixNano(), &types.MarketOrderInfo{
						OrderId:      item.GetUniqueId(),
						IsFinalTrade: true,
					})
				} else {
					e.resultNotify <- e.tradeResult(item, bid, bid.GetPrice(), curTradeQuantity, time.Now().UnixNano(), &types.MarketOrderInfo{
						OrderId:      item.GetUniqueId(),
						IsFinalTrade: false,
					})
				}
				return true
			} else if item.GetSubOrderType() == types.SubOrderTypeMarketByAmount {
				//市价-按成交金额成交
				if bid.GetPrice().Cmp(decimal.Zero) <= 0 {
					return false
				}

				//计算最大需要卖出的数量
				maxQty := func(amount, price, needQty decimal.Decimal) decimal.Decimal {
					a := amount.Div(price).Truncate(e.opts.quantityDecimals)
					return decimal.Min(a, needQty).Truncate(e.opts.quantityDecimals)
				}
				maxTradeQty := maxQty(item.GetAmount(), bid.GetPrice(), item.GetHoldQty())

				curTradeQty := decimal.Zero
				if maxTradeQty.Cmp(e.opts.minTradeQuantity) < 0 {
					return false
				}

				if bid.GetQuantity().Cmp(maxTradeQty) <= 0 {
					curTradeQty = bid.GetQuantity()
					e.bids.Remove(bid.GetUniqueId())
				} else {
					curTradeQty = maxTradeQty
					e.bids.SetQuantity(bid, bid.GetQuantity().Sub(curTradeQty))
				}

				if curTradeQty.Equal(decimal.Zero) {
					return false
				}

				item.SetAmount(item.GetAmount().Sub(curTradeQty.Mul(bid.GetPrice())))
				//市价 按成交额卖出时，需要用户持有的资产数量来进行限制
				item.SetHoldQty(item.GetHoldQty().Sub(curTradeQty))

				//退出条件
				// a.对面订单空了
				// b.金额完全成交
				// c.剩余资金不满足最小成交量
				if e.bids.Len() == 0 ||
					maxQty(item.GetAmount(), e.bids.Top().GetPrice(), item.GetHoldQty()).Cmp(e.opts.minTradeQuantity) <= 0 {
					e.resultNotify <- e.tradeResult(item, bid, bid.GetPrice(), curTradeQty, time.Now().UnixNano(), &types.MarketOrderInfo{
						OrderId:      item.GetUniqueId(),
						IsFinalTrade: true,
					})
				} else {
					e.resultNotify <- e.tradeResult(item, bid, bid.GetPrice(), curTradeQty, time.Now().UnixNano(), &types.MarketOrderInfo{
						OrderId:      item.GetUniqueId(),
						IsFinalTrade: false,
					})
				}

				return true
			}

			return false
		}()

		if !ok {
			//市价单都需要触发一个成交后取消剩余部分的信号
			go func() {
				// time.Sleep(time.Second)
				e.removeNotify <- types.RemoveResult{
					Symbol:   e.symbol,
					UniqueId: item.GetUniqueId(),
					Type:     types.RemoveItemTypeByMarket,
				}
			}()
			break
		} else {
			trade_cnt++
		}

	}
}
