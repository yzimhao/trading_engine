package trading_core

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
)

var btcusdt = NewTradePair("btcusdt", 2, 2)

func d(f float64) decimal.Decimal {
	return decimal.NewFromFloat(f)
}

func TestTradePairDepth(t *testing.T) {

	Convey("卖盘深度行情", t, func() {
		btcusdt.cleanAll()
		btcusdt.PushNewOrder(NewAskLimitItem("id1", d(1.01), d(2), 1112))
		btcusdt.PushNewOrder(NewAskLimitItem("id2", d(1.01), d(2), 1113))
		btcusdt.PushNewOrder(NewAskLimitItem("id3", d(1.1), d(2), 1114))
		time.Sleep(time.Millisecond * time.Duration(100))
		ask := btcusdt.GetAskDepth(0)
		str_ask, _ := json.Marshal(ask)
		So(string(str_ask), ShouldEqual, `[["1.01","4.00"],["1.10","2.00"]]`)
	})

	Convey("买盘深度行情", t, func() {
		btcusdt.cleanAll()
		btcusdt.PushNewOrder(NewBidLimitItem("id4", d(1.02), d(2), 1115))
		btcusdt.PushNewOrder(NewBidLimitItem("id5", d(1.3), d(2), 1116))
		btcusdt.PushNewOrder(NewBidLimitItem("id6", d(1.02), d(2), 1117))
		btcusdt.PushNewOrder(NewBidLimitItem("id7", d(0.02), d(1), 1118))

		time.Sleep(time.Millisecond * time.Duration(100))
		bid := btcusdt.GetBidDepth(0)
		str_bid, _ := json.Marshal(bid)
		So(string(str_bid), ShouldEqual, `[["1.30","2.00"],["1.02","4.00"],["0.02","1.00"]]`)
	})
}

func TestTradeFunc_LimitOrder(t *testing.T) {
	//创建一个买单
	Convey("新增一个限价买单", t, func() {
		btcusdt.cleanAll()
		btcusdt.PushNewOrder(NewBidLimitItem("id11", d(1.1), d(1.2), 1112))

		So(btcusdt.askQueue.Len(), ShouldEqual, 0)
		So(btcusdt.bidQueue.Len(), ShouldEqual, 1)
		So(btcusdt.bidQueue.Top().GetPrice(), ShouldEqual, d(1.1))
		So(btcusdt.bidQueue.Top().GetUniqueId(), ShouldEqual, "id11")
		So(btcusdt.bidQueue.Top().GetQuantity(), ShouldEqual, d(1.2))
		So(btcusdt.bidQueue.Top().GetCreateTime(), ShouldEqual, 1112)
	})

	Convey("新增一个限价卖单", t, func() {
		btcusdt.cleanAll()
		btcusdt.PushNewOrder(NewAskLimitItem("id12", d(1.1), d(1.2), 1112))

		So(btcusdt.askQueue.Len(), ShouldEqual, 1)
		So(btcusdt.bidQueue.Len(), ShouldEqual, 0)
		So(btcusdt.askQueue.Top().GetPrice(), ShouldEqual, d(1.1))
		So(btcusdt.askQueue.Top().GetUniqueId(), ShouldEqual, "id12")
		So(btcusdt.askQueue.Top().GetQuantity(), ShouldEqual, d(1.2))
		So(btcusdt.askQueue.Top().GetCreateTime(), ShouldEqual, 1112)

	})

	Convey("清空买卖队列", t, func() {
		btcusdt.cleanAll()
		So(btcusdt.askQueue.Len(), ShouldEqual, 0)
		So(btcusdt.bidQueue.Len(), ShouldEqual, 0)
	})

	Convey("限价买卖单，价格一致，完全成交", t, func() {
		btcusdt.cleanAll()
		btcusdt.PushNewOrder(NewAskLimitItem("id13", d(1.1), d(1.2), time.Now().UnixNano()))
		btcusdt.PushNewOrder(NewBidLimitItem("id23", d(1.1), d(1.2), time.Now().UnixNano()))
		tradeLog := <-btcusdt.ChTradeResult
		So(tradeLog.AskOrderId, ShouldEqual, "id13")
		So(tradeLog.BidOrderId, ShouldEqual, "id23")
		So(tradeLog.TradePrice, ShouldEqual, d(1.1))
		So(tradeLog.TradeQuantity, ShouldEqual, d(1.2))
		//ask提供流动性，bid主动成交 trade_by=buyer
		So(tradeLog.TradeBy, ShouldEqual, ByBuyer)
	})

	Convey("限价买卖单，价格一致，买单部分成交", t, func() {
		btcusdt.cleanAll()
		btcusdt.PushNewOrder(NewBidLimitItem("id24", d(1.1), d(2.3), time.Now().UnixNano()))
		btcusdt.PushNewOrder(NewAskLimitItem("id14", d(1.1), d(1.2), time.Now().UnixNano()))
		tradeLog := <-btcusdt.ChTradeResult

		So(tradeLog.AskOrderId, ShouldEqual, "id14")
		So(tradeLog.BidOrderId, ShouldEqual, "id24")
		So(tradeLog.TradePrice, ShouldEqual, d(1.1))
		So(tradeLog.TradeQuantity, ShouldEqual, d(1.2))
		So(btcusdt.bidQueue.Top().GetQuantity(), ShouldEqual, d(1.1))
		So(btcusdt.bidQueue.Len(), ShouldEqual, 1)
		So(btcusdt.askQueue.Len(), ShouldEqual, 0)
		//bid提供流动性，ask主动成交 trade_by=seller
		So(tradeLog.TradeBy, ShouldEqual, BySeller)
	})

	Convey("限价买卖单，价格一致，卖单部分成交", t, func() {
		btcusdt.cleanAll()
		btcusdt.PushNewOrder(NewAskLimitItem("id15", d(1.1), d(2.2), 1112))
		btcusdt.PushNewOrder(NewBidLimitItem("id25", d(1.1), d(1.3), 1113))
		tradeLog := <-btcusdt.ChTradeResult

		So(tradeLog.AskOrderId, ShouldEqual, "id15")
		So(tradeLog.BidOrderId, ShouldEqual, "id25")
		So(tradeLog.TradePrice, ShouldEqual, d(1.1))
		So(tradeLog.TradeQuantity, ShouldEqual, d(1.3))
		So(btcusdt.askQueue.Top().GetQuantity(), ShouldEqual, d(0.9))
		So(btcusdt.askQueue.Len(), ShouldEqual, 1)
		So(btcusdt.bidQueue.Len(), ShouldEqual, 0)
	})

	Convey("价格一致，时间优先", t, func() {
		btcusdt.cleanAll()
		btcusdt.PushNewOrder(NewAskLimitItem("id16", d(1.1), d(2.2), 1112))
		btcusdt.PushNewOrder(NewAskLimitItem("id26", d(1.1), d(2.2), 1110))

		btcusdt.PushNewOrder(NewBidLimitItem("id36", d(1.1), d(1.3), 1113))
		tradeLog := <-btcusdt.ChTradeResult

		So(tradeLog.AskOrderId, ShouldEqual, "id26")
		So(tradeLog.BidOrderId, ShouldEqual, "id36")
		So(tradeLog.TradePrice, ShouldEqual, d(1.1))
		So(tradeLog.TradeQuantity, ShouldEqual, d(1.3))
		So(btcusdt.askQueue.Top().GetQuantity(), ShouldEqual, d(0.9))
		So(btcusdt.askQueue.Len(), ShouldEqual, 2)
		So(btcusdt.bidQueue.Len(), ShouldEqual, 0)
	})

	//价格优先
	Convey("价格优先", t, func() {
		btcusdt.cleanAll()
		btcusdt.PushNewOrder(NewAskLimitItem("id17", d(1.01), d(2.2), 1112))
		btcusdt.PushNewOrder(NewAskLimitItem("id27", d(1.1), d(2.2), 1110))

		btcusdt.PushNewOrder(NewBidLimitItem("id37", d(1.1), d(1.3), 1113))
		tradeLog := <-btcusdt.ChTradeResult

		So(tradeLog.AskOrderId, ShouldEqual, "id17")
		So(tradeLog.BidOrderId, ShouldEqual, "id37")
		So(tradeLog.TradePrice, ShouldEqual, d(1.01))
		So(tradeLog.TradeQuantity, ShouldEqual, d(1.3))
		So(tradeLog.Last, ShouldEqual, "")
		So(btcusdt.askQueue.Top().GetQuantity(), ShouldEqual, d(0.9))
		So(btcusdt.askQueue.Len(), ShouldEqual, 2)
		So(btcusdt.bidQueue.Len(), ShouldEqual, 0)
	})

}

func TestTradeFunc_MarketBuyOrder(t *testing.T) {

	Convey("市价买入 按数量买入, 金额足够买单完全成交", t, func() {
		btcusdt.cleanAll()
		btcusdt.PushNewOrder(NewAskLimitItem("id18", d(1.01), d(2.2), 1112))
		btcusdt.PushNewOrder(NewBidMarketQtyItem("id28", d(1.1), d(100), 1113))
		tradeLog := <-btcusdt.ChTradeResult
		<-btcusdt.ChCancelResult

		So(tradeLog.AskOrderId, ShouldEqual, "id18")
		So(tradeLog.BidOrderId, ShouldEqual, "id28")
		So(tradeLog.TradePrice, ShouldEqual, d(1.01))
		So(tradeLog.TradeQuantity, ShouldEqual, d(1.1))
		So(tradeLog.Last, ShouldEqual, "id28")

	})

	Convey("市价按数量买入, 账户金额足够买单部分成交", t, func() {
		btcusdt.cleanAll()
		btcusdt.PushNewOrder(NewAskLimitItem("id19", d(1.01), d(2.2), 1112))
		btcusdt.PushNewOrder(NewBidMarketQtyItem("id29", d(100), d(100), 1113))
		tradeLog := <-btcusdt.ChTradeResult
		<-btcusdt.ChCancelResult

		So(tradeLog.AskOrderId, ShouldEqual, "id19")
		So(tradeLog.BidOrderId, ShouldEqual, "id29")
		So(tradeLog.TradePrice, ShouldEqual, d(1.01))
		So(tradeLog.TradeQuantity, ShouldEqual, d(2.2))
		So(tradeLog.Last, ShouldEqual, "id29")
	})

	//市价买入 按数量, 金额不足 买单部分成交
	Convey("市价按数量买入, 账户金额不足, 买单部分成交", t, func() {
		btcusdt.cleanAll()
		btcusdt.PushNewOrder(NewAskLimitItem("id110", d(100), d(20), 1112))
		btcusdt.PushNewOrder(NewBidMarketQtyItem("id210", d(20), d(100), 1113))
		tradeLog := <-btcusdt.ChTradeResult
		<-btcusdt.ChCancelResult

		So(tradeLog.AskOrderId, ShouldEqual, "id110")
		So(tradeLog.BidOrderId, ShouldEqual, "id210")
		So(tradeLog.TradePrice, ShouldEqual, d(100))
		So(tradeLog.TradeQuantity, ShouldEqual, d(1))

		So(tradeLog.Last, ShouldEqual, "id210")
		So(btcusdt.askQueue.Len(), ShouldEqual, 1)
		So(btcusdt.bidQueue.Len(), ShouldEqual, 0)
	})

	Convey("市价买入 指定金额, 买单完全成交", t, func() {
		btcusdt.cleanAll()
		btcusdt.PushNewOrder(NewAskLimitItem("id111", d(10.00), d(100), 1112))
		btcusdt.PushNewOrder(NewBidMarketAmountItem("id211", d(50), 1113))
		tradeLog := <-btcusdt.ChTradeResult
		<-btcusdt.ChCancelResult

		So(tradeLog.AskOrderId, ShouldEqual, "id111")
		So(tradeLog.BidOrderId, ShouldEqual, "id211")
		So(tradeLog.TradePrice, ShouldEqual, d(10.00))
		So(tradeLog.TradeQuantity, ShouldEqual, d(5))

		So(tradeLog.Last, ShouldEqual, "id211")
		So(btcusdt.askQueue.Len(), ShouldEqual, 1)
	})

	Convey("市价买入 指定金额, 买单部分成交", t, func() {
		btcusdt.cleanAll()
		btcusdt.PushNewOrder(NewAskLimitItem("id112", d(10.00), d(100), 1112))
		btcusdt.PushNewOrder(NewBidMarketAmountItem("id212", d(6000), 1113))
		tradeLog := <-btcusdt.ChTradeResult
		<-btcusdt.ChCancelResult

		So(tradeLog.AskOrderId, ShouldEqual, "id112")
		So(tradeLog.BidOrderId, ShouldEqual, "id212")
		So(tradeLog.TradePrice, ShouldEqual, d(10.00))
		So(tradeLog.TradeQuantity, ShouldEqual, d(100))

		So(tradeLog.Last, ShouldEqual, "id212")
		So(btcusdt.askQueue.Len(), ShouldEqual, 0)
		So(btcusdt.bidQueue.Len(), ShouldEqual, 0)
	})
}

func TestTradeFunc_MarketSellOrder(t *testing.T) {
	Convey("市价卖出 按数量, 完全成交", t, func() {

		btcusdt.cleanAll()
		btcusdt.PushNewOrder(NewBidLimitItem("id113", d(10.00), d(100), 1112))
		btcusdt.PushNewOrder(NewAskMarketQtyItem("id213", d(6), 1113))
		tradeLog := <-btcusdt.ChTradeResult
		<-btcusdt.ChCancelResult

		So(tradeLog.AskOrderId, ShouldEqual, "id213")
		So(tradeLog.BidOrderId, ShouldEqual, "id113")
		So(tradeLog.TradePrice, ShouldEqual, d(10.00))
		So(tradeLog.TradeQuantity, ShouldEqual, d(6))

		So(tradeLog.Last, ShouldEqual, "id213")
		So(btcusdt.askQueue.Len(), ShouldEqual, 0)
		So(btcusdt.bidQueue.Len(), ShouldEqual, 1)
	})

	Convey("市价卖出 按数量，部分成交", t, func() {
		//
		btcusdt.cleanAll()
		btcusdt.PushNewOrder(NewBidLimitItem("id114", d(10.00), d(100), 1112))
		btcusdt.PushNewOrder(NewAskMarketQtyItem("id214", d(6000), 1113))
		tradeLog := <-btcusdt.ChTradeResult
		<-btcusdt.ChCancelResult

		So(tradeLog.AskOrderId, ShouldEqual, "id214")
		So(tradeLog.BidOrderId, ShouldEqual, "id114")
		So(tradeLog.TradePrice, ShouldEqual, d(10.00))
		So(tradeLog.TradeQuantity, ShouldEqual, d(100))

		So(tradeLog.Last, ShouldEqual, "id214")
		So(btcusdt.askQueue.Len(), ShouldEqual, 0)
		So(btcusdt.bidQueue.Len(), ShouldEqual, 0)
	})

	Convey("市价卖出 指定金额，持仓足够完全成交", t, func() {
		btcusdt.cleanAll()
		btcusdt.PushNewOrder(NewBidLimitItem("id115", d(10.00), d(1000), 1112))
		btcusdt.PushNewOrder(NewAskMarketAmountItem("id215", d(6000), d(1000000), 1113))
		tradeLog := <-btcusdt.ChTradeResult
		<-btcusdt.ChCancelResult

		So(tradeLog.AskOrderId, ShouldEqual, "id215")
		So(tradeLog.BidOrderId, ShouldEqual, "id115")
		So(tradeLog.TradePrice, ShouldEqual, d(10.00))
		So(tradeLog.TradeQuantity, ShouldEqual, d(600))

		So(tradeLog.Last, ShouldEqual, "id215")
		So(btcusdt.askQueue.Len(), ShouldEqual, 0)
		So(btcusdt.bidQueue.Len(), ShouldEqual, 1)
	})

	Convey("市价卖出 指定金额，持仓足够 但是部分成交", t, func() {
		btcusdt.cleanAll()
		btcusdt.PushNewOrder(NewBidLimitItem("id116", d(10.00), d(50), 1112))
		btcusdt.PushNewOrder(NewAskMarketAmountItem("id216", d(6000), d(1000000), 1113))
		tradeLog := <-btcusdt.ChTradeResult
		<-btcusdt.ChCancelResult

		So(tradeLog.AskOrderId, ShouldEqual, "id216")
		So(tradeLog.BidOrderId, ShouldEqual, "id116")
		So(tradeLog.TradePrice, ShouldEqual, d(10.00))
		So(tradeLog.TradeQuantity, ShouldEqual, d(50))

		So(tradeLog.Last, ShouldEqual, "id216")
		So(btcusdt.askQueue.Len(), ShouldEqual, 0)
		So(btcusdt.bidQueue.Len(), ShouldEqual, 0)
	})

	Convey("市价卖出 指定金额，持仓不足", t, func() {
		btcusdt.cleanAll()
		btcusdt.PushNewOrder(NewBidLimitItem("id117", d(100.00), d(50), 1112))
		btcusdt.PushNewOrder(NewAskMarketAmountItem("id217", d(500), d(3), 1113))
		tradeLog := <-btcusdt.ChTradeResult
		<-btcusdt.ChCancelResult

		So(tradeLog.AskOrderId, ShouldEqual, "id217")
		So(tradeLog.BidOrderId, ShouldEqual, "id117")
		So(tradeLog.TradePrice, ShouldEqual, d(100))
		So(tradeLog.TradeQuantity, ShouldEqual, d(3))

		So(tradeLog.Last, ShouldEqual, "id217")
		So(btcusdt.askQueue.Len(), ShouldEqual, 0)
		So(btcusdt.bidQueue.Len(), ShouldEqual, 1)
	})

	Convey("市价卖出 指定金额，一个都没有成交，只输出市价撤单信号", t, func() {
		btcusdt.cleanAll()

		cancel := make(chan CancelBody)
		go func() {
			for {
				select {
				case cancel <- <-btcusdt.ChCancelResult:

				}
			}
		}()

		btcusdt.PushNewOrder(NewBidLimitItem("id1118", d(1000.00), d(50), 1112))
		btcusdt.PushNewOrder(NewAskMarketAmountItem("id2218", d(1), d(30), 1113))

		So(btcusdt.bidQueue.Len(), ShouldEqual, 1)
		bid := btcusdt.bidQueue.Top()
		So(bid.GetQuantity(), ShouldEqual, d(50))
		So(btcusdt.askQueue.Len(), ShouldEqual, 0)

		cinfo := <-cancel
		So(cinfo.OrderId, ShouldEqual, "id2218")
	})
}
