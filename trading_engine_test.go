package trading_engine

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var btcusdt = NewTradePair("btcusdt", 2, 0)

func init() {

}

func d(f float64) decimal.Decimal {
	return decimal.NewFromFloat(f)
}

func TestTradePairDepth(t *testing.T) {
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewAskLimitItem("id1", d(1.01), d(2), 1112))
	btcusdt.PushNewOrder(NewAskLimitItem("id2", d(1.01), d(2), 1113))
	btcusdt.PushNewOrder(NewAskLimitItem("id3", d(1.1), d(2), 1114))

	btcusdt.PushNewOrder(NewBidLimitItem("id4", d(1.02), d(2), 1115))
	btcusdt.PushNewOrder(NewBidLimitItem("id5", d(1.3), d(2), 1116))
	btcusdt.PushNewOrder(NewBidLimitItem("id6", d(1.02), d(2), 1117))
	btcusdt.PushNewOrder(NewBidLimitItem("id7", d(0.02), d(1), 1118))

	time.Sleep(time.Millisecond * time.Duration(100))
	ask := btcusdt.GetAskDepth(0)
	str_ask, _ := json.Marshal(ask)
	assert.Equal(t, "[[\"1.01\",\"4\"],[\"1.10\",\"2\"]]", string(str_ask))

	bid := btcusdt.GetBidDepth(0)
	str_bid, _ := json.Marshal(bid)
	assert.Equal(t, "[[\"1.30\",\"2\"],[\"1.02\",\"4\"],[\"0.02\",\"1\"]]", string(str_bid))
}

func TestTradeFunc_LimitOrder(t *testing.T) {
	//创建一个买单
	btcusdt.PushNewOrder(NewBidLimitItem("id1", d(1.1), d(1.2), 1112))
	assert.Equal(t, 0, btcusdt.askQueue.Len())
	assert.Equal(t, 1, btcusdt.bidQueue.Len())
	assert.Equal(t, "id1", btcusdt.bidQueue.Top().GetUniqueId())

	//清空队列
	btcusdt.cleanAll()
	assert.Equal(t, 0, btcusdt.askQueue.Len())
	assert.Equal(t, 0, btcusdt.bidQueue.Len())

	//创建一组买卖单，价格一致，完全成交
	btcusdt.PushNewOrder(NewAskLimitItem("id1", d(1.1), d(1.2), 1112))
	btcusdt.PushNewOrder(NewBidLimitItem("id2", d(1.1), d(1.2), 1113))
	assert.Equal(t, 1, btcusdt.askQueue.Len())
	assert.Equal(t, "id1", btcusdt.askQueue.Top().GetUniqueId())
	tradeLog := <-btcusdt.ChTradeResult
	assert.Equal(t, "id1", tradeLog.AskOrderId)
	assert.Equal(t, "id2", tradeLog.BidOrderId)
	assert.Equal(t, "1.1", tradeLog.TradePrice.String())
	assert.Equal(t, "1.2", tradeLog.TradeQuantity.String())

	//一组订单，价格一致，买单部分成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewAskLimitItem("id1", d(1.1), d(1.2), 1112))
	btcusdt.PushNewOrder(NewBidLimitItem("id2", d(1.1), d(2.3), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "id1", tradeLog.AskOrderId)
	assert.Equal(t, "id2", tradeLog.BidOrderId)
	assert.Equal(t, "1.1", tradeLog.TradePrice.String())
	assert.Equal(t, "1.2", tradeLog.TradeQuantity.String())
	assert.Equal(t, "1.1", btcusdt.bidQueue.Top().GetQuantity().String())
	assert.Equal(t, 1, btcusdt.bidQueue.Len())
	assert.Equal(t, 0, btcusdt.askQueue.Len())

	//一组订单，价格一致，卖单部分成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewAskLimitItem("id1", d(1.1), d(2.2), 1112))
	btcusdt.PushNewOrder(NewBidLimitItem("id2", d(1.1), d(1.3), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "id1", tradeLog.AskOrderId)
	assert.Equal(t, "id2", tradeLog.BidOrderId)
	assert.Equal(t, "1.1", tradeLog.TradePrice.String())
	assert.Equal(t, "1.3", tradeLog.TradeQuantity.String())
	assert.Equal(t, "0.9", btcusdt.askQueue.Top().GetQuantity().String())
	assert.Equal(t, 0, btcusdt.bidQueue.Len())
	assert.Equal(t, 1, btcusdt.askQueue.Len())

	//时间优先
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewAskLimitItem("id1", d(1.1), d(2.2), 1112))
	btcusdt.PushNewOrder(NewAskLimitItem("id2", d(1.1), d(2.2), 1110))

	btcusdt.PushNewOrder(NewBidLimitItem("id3", d(1.1), d(1.3), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "id2", tradeLog.AskOrderId)
	assert.Equal(t, "id3", tradeLog.BidOrderId)
	assert.Equal(t, "1.1", tradeLog.TradePrice.String())
	assert.Equal(t, "1.3", tradeLog.TradeQuantity.String())
	assert.Equal(t, "0.9", btcusdt.askQueue.Top().GetQuantity().String())
	assert.Equal(t, 0, btcusdt.bidQueue.Len())
	assert.Equal(t, 2, btcusdt.askQueue.Len())

	//价格优先
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewAskLimitItem("id1", d(1.01), d(2.2), 1112))
	btcusdt.PushNewOrder(NewAskLimitItem("id2", d(1.1), d(2.2), 1110))

	btcusdt.PushNewOrder(NewBidLimitItem("id3", d(1.1), d(1.3), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "id1", tradeLog.AskOrderId)
	assert.Equal(t, "id3", tradeLog.BidOrderId)
	assert.Equal(t, "1.01", tradeLog.TradePrice.String())
	assert.Equal(t, "1.3", tradeLog.TradeQuantity.String())
	assert.Equal(t, "0.9", btcusdt.askQueue.Top().GetQuantity().String())
	assert.Equal(t, 0, btcusdt.bidQueue.Len())
	assert.Equal(t, 2, btcusdt.askQueue.Len())

}

func TestTradeFunc_MarketBuyOrder(t *testing.T) {

	//市价买入 按数量, 金额足够买单完全成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewAskLimitItem("id1", d(1.01), d(2.2), 1112))
	btcusdt.PushNewOrder(NewBidMarketQtyItem("id2", d(1.1), d(100), 1113))
	tradeLog := <-btcusdt.ChTradeResult
	assert.Equal(t, "id1", tradeLog.AskOrderId)
	assert.Equal(t, "id2", tradeLog.BidOrderId)
	assert.Equal(t, "1.01", tradeLog.TradePrice.String())
	assert.Equal(t, "1.1", tradeLog.TradeQuantity.String())

	//市价买入 按数量, 金额足够买单部分成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewAskLimitItem("id1", d(1.01), d(2.2), 1112))
	btcusdt.PushNewOrder(NewBidMarketQtyItem("id2", d(100), d(100), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "id1", tradeLog.AskOrderId)
	assert.Equal(t, "id2", tradeLog.BidOrderId)
	assert.Equal(t, "1.01", tradeLog.TradePrice.String())
	assert.Equal(t, "2.2", tradeLog.TradeQuantity.String())
	cancelOrder := <-btcusdt.ChCancelResult
	assert.Equal(t, "id2", cancelOrder)
	assert.Equal(t, 0, btcusdt.askQueue.Len())

	//市价买入 按数量, 金额不足 买单部分成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewAskLimitItem("id1", d(100), d(20), 1112))
	btcusdt.PushNewOrder(NewBidMarketQtyItem("id2", d(20), d(100), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "id1", tradeLog.AskOrderId)
	assert.Equal(t, "id2", tradeLog.BidOrderId)
	assert.Equal(t, "100", tradeLog.TradePrice.String())
	assert.Equal(t, "1", tradeLog.TradeQuantity.String())
	cancelOrder = <-btcusdt.ChCancelResult
	assert.Equal(t, "id2", cancelOrder)
	assert.Equal(t, 1, btcusdt.askQueue.Len())
	assert.Equal(t, 0, btcusdt.bidQueue.Len())

	//市价买入 指定金额, 买单完全成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewAskLimitItem("id1", d(10.00), d(100), 1112))
	btcusdt.PushNewOrder(NewBidMarketAmountItem("id2", d(50), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "id1", tradeLog.AskOrderId)
	assert.Equal(t, "id2", tradeLog.BidOrderId)
	assert.Equal(t, "10", tradeLog.TradePrice.String())
	assert.Equal(t, "5", tradeLog.TradeQuantity.String())
	cancelOrder = <-btcusdt.ChCancelResult
	assert.Equal(t, "id2", cancelOrder)
	assert.Equal(t, 1, btcusdt.askQueue.Len())

	//市价买入 指定金额, 买单部分成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewAskLimitItem("id1", d(10.00), d(100), 1112))
	btcusdt.PushNewOrder(NewBidMarketAmountItem("id2", d(6000), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "id1", tradeLog.AskOrderId)
	assert.Equal(t, "id2", tradeLog.BidOrderId)
	assert.Equal(t, "10", tradeLog.TradePrice.String())
	assert.Equal(t, "100", tradeLog.TradeQuantity.String())
	cancelOrder = <-btcusdt.ChCancelResult
	assert.Equal(t, "id2", cancelOrder)
	assert.Equal(t, 0, btcusdt.askQueue.Len())
	assert.Equal(t, 0, btcusdt.bidQueue.Len())

}

func TestTradeFunc_MarketSellOrder(t *testing.T) {
	//市价卖出 按数量, 完全成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewBidLimitItem("id1", d(10.00), d(100), 1112))
	btcusdt.PushNewOrder(NewAskMarketQtyItem("id2", d(6), 1113))
	tradeLog := <-btcusdt.ChTradeResult
	assert.Equal(t, "id2", tradeLog.AskOrderId)
	assert.Equal(t, "id1", tradeLog.BidOrderId)
	assert.Equal(t, "10", tradeLog.TradePrice.String())
	assert.Equal(t, "6", tradeLog.TradeQuantity.String())
	cancelOrder := <-btcusdt.ChCancelResult
	assert.Equal(t, "id2", cancelOrder)
	assert.Equal(t, 0, btcusdt.askQueue.Len())
	assert.Equal(t, 1, btcusdt.bidQueue.Len())

	//市价卖出 按数量，部分成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewBidLimitItem("id1", d(10.00), d(100), 1112))
	btcusdt.PushNewOrder(NewAskMarketQtyItem("id2", d(6000), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "id2", tradeLog.AskOrderId)
	assert.Equal(t, "id1", tradeLog.BidOrderId)
	assert.Equal(t, "10", tradeLog.TradePrice.String())
	assert.Equal(t, "100", tradeLog.TradeQuantity.String())
	cancelOrder = <-btcusdt.ChCancelResult
	assert.Equal(t, "id2", cancelOrder)
	assert.Equal(t, 0, btcusdt.askQueue.Len())
	assert.Equal(t, 0, btcusdt.bidQueue.Len())

	//市价卖出 指定金额，持仓足够完全成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewBidLimitItem("id1", d(10.00), d(1000), 1112))
	btcusdt.PushNewOrder(NewAskMarketAmountItem("id2", d(6000), d(1000000), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "id2", tradeLog.AskOrderId)
	assert.Equal(t, "id1", tradeLog.BidOrderId)
	assert.Equal(t, "10", tradeLog.TradePrice.String())
	assert.Equal(t, "600", tradeLog.TradeQuantity.String())
	cancelOrder = <-btcusdt.ChCancelResult
	assert.Equal(t, "id2", cancelOrder)
	assert.Equal(t, 0, btcusdt.askQueue.Len())
	assert.Equal(t, 1, btcusdt.bidQueue.Len())

	//市价卖出 指定金额，持仓足够 但是部分成交
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewBidLimitItem("id1", d(10.00), d(50), 1112))
	btcusdt.PushNewOrder(NewAskMarketAmountItem("id2", d(6000), d(1000000), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "id2", tradeLog.AskOrderId)
	assert.Equal(t, "id1", tradeLog.BidOrderId)
	assert.Equal(t, "10", tradeLog.TradePrice.String())
	assert.Equal(t, "50", tradeLog.TradeQuantity.String())
	cancelOrder = <-btcusdt.ChCancelResult
	assert.Equal(t, "id2", cancelOrder)
	assert.Equal(t, 0, btcusdt.askQueue.Len())
	assert.Equal(t, 0, btcusdt.bidQueue.Len())

	//市价卖出 指定金额，持仓不足
	btcusdt.cleanAll()
	btcusdt.PushNewOrder(NewBidLimitItem("id1", d(100.00), d(50), 1112))
	btcusdt.PushNewOrder(NewAskMarketAmountItem("id2", d(500), d(3), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "id2", tradeLog.AskOrderId)
	assert.Equal(t, "id1", tradeLog.BidOrderId)
	assert.Equal(t, "100", tradeLog.TradePrice.String())
	assert.Equal(t, "3", tradeLog.TradeQuantity.String())
	cancelOrder = <-btcusdt.ChCancelResult
	assert.Equal(t, "id2", cancelOrder)
	assert.Equal(t, 0, btcusdt.askQueue.Len())
	assert.Equal(t, 1, btcusdt.bidQueue.Len())
}
