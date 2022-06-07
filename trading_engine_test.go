package trading_engine

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var btcusdt = NewTradePair("btcusdt", 2, 0)

func init() {

}

func d(f float64) decimal.Decimal {
	return decimal.NewFromFloat(f)
}

func TestTradeFunc(t *testing.T) {

}

func TestTradeFunc_LimitOrder(t *testing.T) {
	//创建一个买单
	btcusdt.PushNewOrder(NewBidLimitItem("uid1", d(1.1), d(1.2), 1112))
	assert.Equal(t, 0, btcusdt.askQueue.Len())
	assert.Equal(t, 1, btcusdt.bidQueue.Len())
	assert.Equal(t, "uid1", btcusdt.bidQueue.Top().GetUniqueId())

	//清空队列
	btcusdt.cleanAll()
	assert.Equal(t, 0, btcusdt.askQueue.Len())
	assert.Equal(t, 0, btcusdt.bidQueue.Len())

	// time.Sleep(time.Duration(100) * time.Millisecond)
	// a := btcusdt.GetBidepth(0)
	// stra, _ := json.Marshal(a)
	// assert.Equal(t, "[[\"1.10\",\"1\"]]", string(stra))

	//创建一组买卖单，价格一致，完全成交
	btcusdt.PushNewOrder(NewAskLimitItem("uid1", d(1.1), d(1.2), 1112))
	btcusdt.PushNewOrder(NewBidLimitItem("uid2", d(1.1), d(1.2), 1113))
	assert.Equal(t, 1, btcusdt.askQueue.Len())
	assert.Equal(t, "uid1", btcusdt.askQueue.Top().GetUniqueId())
	tradeLog := <-btcusdt.ChTradeResult
	assert.Equal(t, "uid1", tradeLog.AskOrderId)
	assert.Equal(t, "uid2", tradeLog.BidOrderId)
	assert.Equal(t, "1.1", tradeLog.TradePrice.String())
	assert.Equal(t, "1.2", tradeLog.TradeQuantity.String())

	//一组订单，价格一致，买单部分成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewAskLimitItem("uid1", d(1.1), d(1.2), 1112))
	btcusdt.PushNewOrder(NewBidLimitItem("uid2", d(1.1), d(2.3), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "uid1", tradeLog.AskOrderId)
	assert.Equal(t, "uid2", tradeLog.BidOrderId)
	assert.Equal(t, "1.1", tradeLog.TradePrice.String())
	assert.Equal(t, "1.2", tradeLog.TradeQuantity.String())
	assert.Equal(t, "1.1", btcusdt.bidQueue.Top().GetQuantity().String())
	assert.Equal(t, 1, btcusdt.bidQueue.Len())
	assert.Equal(t, 0, btcusdt.askQueue.Len())

	//一组订单，价格一致，卖单部分成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewAskLimitItem("uid1", d(1.1), d(2.2), 1112))
	btcusdt.PushNewOrder(NewBidLimitItem("uid2", d(1.1), d(1.3), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "uid1", tradeLog.AskOrderId)
	assert.Equal(t, "uid2", tradeLog.BidOrderId)
	assert.Equal(t, "1.1", tradeLog.TradePrice.String())
	assert.Equal(t, "1.3", tradeLog.TradeQuantity.String())
	assert.Equal(t, "0.9", btcusdt.askQueue.Top().GetQuantity().String())
	assert.Equal(t, 0, btcusdt.bidQueue.Len())
	assert.Equal(t, 1, btcusdt.askQueue.Len())

	//时间优先
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewAskLimitItem("uid1", d(1.1), d(2.2), 1112))
	btcusdt.PushNewOrder(NewAskLimitItem("uid2", d(1.1), d(2.2), 1110))

	btcusdt.PushNewOrder(NewBidLimitItem("uid3", d(1.1), d(1.3), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "uid2", tradeLog.AskOrderId)
	assert.Equal(t, "uid3", tradeLog.BidOrderId)
	assert.Equal(t, "1.1", tradeLog.TradePrice.String())
	assert.Equal(t, "1.3", tradeLog.TradeQuantity.String())
	assert.Equal(t, "0.9", btcusdt.askQueue.Top().GetQuantity().String())
	assert.Equal(t, 0, btcusdt.bidQueue.Len())
	assert.Equal(t, 2, btcusdt.askQueue.Len())

	//价格优先
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewAskLimitItem("uid1", d(1.01), d(2.2), 1112))
	btcusdt.PushNewOrder(NewAskLimitItem("uid2", d(1.1), d(2.2), 1110))

	btcusdt.PushNewOrder(NewBidLimitItem("uid3", d(1.1), d(1.3), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "uid1", tradeLog.AskOrderId)
	assert.Equal(t, "uid3", tradeLog.BidOrderId)
	assert.Equal(t, "1.01", tradeLog.TradePrice.String())
	assert.Equal(t, "1.3", tradeLog.TradeQuantity.String())
	assert.Equal(t, "0.9", btcusdt.askQueue.Top().GetQuantity().String())
	assert.Equal(t, 0, btcusdt.bidQueue.Len())
	assert.Equal(t, 2, btcusdt.askQueue.Len())

}

func TestTradeFunc_MarketBuyOrder(t *testing.T) {

	//市价买入 按数量, 金额足够买单完全成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewAskLimitItem("uid1", d(1.01), d(2.2), 1112))
	btcusdt.PushNewOrder(NewBidMarketQtyItem("uid2", d(1.1), d(100), 1113))
	tradeLog := <-btcusdt.ChTradeResult
	assert.Equal(t, "uid1", tradeLog.AskOrderId)
	assert.Equal(t, "uid2", tradeLog.BidOrderId)
	assert.Equal(t, "1.01", tradeLog.TradePrice.String())
	assert.Equal(t, "1.1", tradeLog.TradeQuantity.String())

	//市价买入 按数量, 金额足够买单部分成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewAskLimitItem("uid1", d(1.01), d(2.2), 1112))
	btcusdt.PushNewOrder(NewBidMarketQtyItem("uid2", d(100), d(100), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "uid1", tradeLog.AskOrderId)
	assert.Equal(t, "uid2", tradeLog.BidOrderId)
	assert.Equal(t, "1.01", tradeLog.TradePrice.String())
	assert.Equal(t, "2.2", tradeLog.TradeQuantity.String())
	cancelOrder := <-btcusdt.ChCancelResult
	assert.Equal(t, "uid2", cancelOrder)
	assert.Equal(t, 0, btcusdt.askQueue.Len())

	//市价买入 按数量, 金额不足 买单部分成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewAskLimitItem("uid1", d(100), d(20), 1112))
	btcusdt.PushNewOrder(NewBidMarketQtyItem("uid2", d(20), d(100), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "uid1", tradeLog.AskOrderId)
	assert.Equal(t, "uid2", tradeLog.BidOrderId)
	assert.Equal(t, "100", tradeLog.TradePrice.String())
	assert.Equal(t, "5", tradeLog.TradeQuantity.String())
	cancelOrder = <-btcusdt.ChCancelResult
	assert.Equal(t, "uid2", cancelOrder)
	assert.Equal(t, 1, btcusdt.askQueue.Len())
	assert.Equal(t, 0, btcusdt.bidQueue.Len())

	//市价买入 指定金额, 买单完全成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewAskLimitItem("uid1", d(10.00), d(100), 1112))
	btcusdt.PushNewOrder(NewBidMarketAmountItem("uid2", d(50), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "uid1", tradeLog.AskOrderId)
	assert.Equal(t, "uid2", tradeLog.BidOrderId)
	assert.Equal(t, "10", tradeLog.TradePrice.String())
	assert.Equal(t, "5", tradeLog.TradeQuantity.String())
	cancelOrder = <-btcusdt.ChCancelResult
	assert.Equal(t, "uid2", cancelOrder)
	assert.Equal(t, 1, btcusdt.askQueue.Len())

	//市价买入 指定金额, 买单部分成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewAskLimitItem("uid1", d(10.00), d(100), 1112))
	btcusdt.PushNewOrder(NewBidMarketAmountItem("uid2", d(6000), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "uid1", tradeLog.AskOrderId)
	assert.Equal(t, "uid2", tradeLog.BidOrderId)
	assert.Equal(t, "10", tradeLog.TradePrice.String())
	assert.Equal(t, "100", tradeLog.TradeQuantity.String())
	cancelOrder = <-btcusdt.ChCancelResult
	assert.Equal(t, "uid2", cancelOrder)
	assert.Equal(t, 0, btcusdt.askQueue.Len())
	assert.Equal(t, 0, btcusdt.bidQueue.Len())

}

func TestTradeFunc_MarketSellOrder(t *testing.T) {
	//市价卖出 按数量, 完全成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewBidLimitItem("uid1", d(10.00), d(100), 1112))
	btcusdt.PushNewOrder(NewAskMarketQtyItem("uid2", d(6), 1113))
	tradeLog := <-btcusdt.ChTradeResult
	assert.Equal(t, "uid2", tradeLog.AskOrderId)
	assert.Equal(t, "uid1", tradeLog.BidOrderId)
	assert.Equal(t, "10", tradeLog.TradePrice.String())
	assert.Equal(t, "6", tradeLog.TradeQuantity.String())
	cancelOrder := <-btcusdt.ChCancelResult
	assert.Equal(t, "uid2", cancelOrder)
	assert.Equal(t, 0, btcusdt.askQueue.Len())
	assert.Equal(t, 1, btcusdt.bidQueue.Len())

	//市价卖出 按数量，部分成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewBidLimitItem("uid1", d(10.00), d(100), 1112))
	btcusdt.PushNewOrder(NewAskMarketQtyItem("uid2", d(6000), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "uid2", tradeLog.AskOrderId)
	assert.Equal(t, "uid1", tradeLog.BidOrderId)
	assert.Equal(t, "10", tradeLog.TradePrice.String())
	assert.Equal(t, "100", tradeLog.TradeQuantity.String())
	cancelOrder = <-btcusdt.ChCancelResult
	assert.Equal(t, "uid2", cancelOrder)
	assert.Equal(t, 0, btcusdt.askQueue.Len())
	assert.Equal(t, 0, btcusdt.bidQueue.Len())

	//市价卖出 指定金额，持仓足够完全成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewBidLimitItem("uid1", d(10.00), d(1000), 1112))
	btcusdt.PushNewOrder(NewAskMarketAmountItem("uid2", d(6000), d(1000000), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "uid2", tradeLog.AskOrderId)
	assert.Equal(t, "uid1", tradeLog.BidOrderId)
	assert.Equal(t, "10", tradeLog.TradePrice.String())
	assert.Equal(t, "600", tradeLog.TradeQuantity.String())
	cancelOrder = <-btcusdt.ChCancelResult
	assert.Equal(t, "uid2", cancelOrder)
	assert.Equal(t, 0, btcusdt.askQueue.Len())
	assert.Equal(t, 1, btcusdt.bidQueue.Len())

	//市价卖出 指定金额，持仓足够 但是部分成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewBidLimitItem("uid1", d(10.00), d(50), 1112))
	btcusdt.PushNewOrder(NewAskMarketAmountItem("uid2", d(6000), d(1000000), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "uid2", tradeLog.AskOrderId)
	assert.Equal(t, "uid1", tradeLog.BidOrderId)
	assert.Equal(t, "10", tradeLog.TradePrice.String())
	assert.Equal(t, "50", tradeLog.TradeQuantity.String())
	cancelOrder = <-btcusdt.ChCancelResult
	assert.Equal(t, "uid2", cancelOrder)
	assert.Equal(t, 0, btcusdt.askQueue.Len())
	assert.Equal(t, 0, btcusdt.bidQueue.Len())

	//市价卖出 指定金额，持仓不足
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewBidLimitItem("uid1", d(100.00), d(50), 1112))
	btcusdt.PushNewOrder(NewAskMarketAmountItem("uid2", d(500), d(3), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "uid2", tradeLog.AskOrderId)
	assert.Equal(t, "uid1", tradeLog.BidOrderId)
	assert.Equal(t, "100", tradeLog.TradePrice.String())
	assert.Equal(t, "3", tradeLog.TradeQuantity.String())
	cancelOrder = <-btcusdt.ChCancelResult
	assert.Equal(t, "uid2", cancelOrder)
	assert.Equal(t, 0, btcusdt.askQueue.Len())
	assert.Equal(t, 1, btcusdt.bidQueue.Len())
}
