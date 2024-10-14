package matching_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/yzimhao/trading_engine/v2/pkg/matching"
	"github.com/yzimhao/trading_engine/v2/pkg/matching/types"
)

var (
	btcusdt *matching.Engine
)

func init() {
	ctx := context.Background()

	opts := []matching.Option{
		matching.WithPriceDecimals(2),
		matching.WithQuantityDecimals(2),
		matching.WithDebug(true),
	}
	btcusdt = matching.NewEngine(ctx, "btcusdt", opts...)
}

func TestTradePairDepth(t *testing.T) {

	Convey("卖盘深度行情", t, func() {
		btcusdt.Clean()
		btcusdt.AddItem(matching.NewAskLimitItem("id1", d(1.01), d(2), 1112))
		btcusdt.AddItem(matching.NewAskLimitItem("id2", d(1.01), d(2), 1113))
		btcusdt.AddItem(matching.NewAskLimitItem("id3", d(1.1), d(2), 1114))
		time.Sleep(time.Millisecond * time.Duration(100))
		ask := btcusdt.GetAskOrderBook(0)
		str_ask, _ := json.Marshal(ask)
		So(string(str_ask), ShouldEqual, `[["1.01","4.00"],["1.10","2.00"]]`)
	})

	Convey("买盘深度行情", t, func() {
		btcusdt.Clean()
		btcusdt.AddItem(matching.NewBidLimitItem("id4", d(1.02), d(2), 1115))
		btcusdt.AddItem(matching.NewBidLimitItem("id5", d(1.3), d(2), 1116))
		btcusdt.AddItem(matching.NewBidLimitItem("id6", d(1.02), d(2), 1117))
		btcusdt.AddItem(matching.NewBidLimitItem("id7", d(0.02), d(1), 1118))

		time.Sleep(time.Millisecond * time.Duration(100))
		bid := btcusdt.GetBidOrderBook(0)
		str_bid, _ := json.Marshal(bid)
		So(string(str_bid), ShouldEqual, `[["1.30","2.00"],["1.02","4.00"],["0.02","1.00"]]`)
	})
}

func TestTradeFunc_LimitOrder(t *testing.T) {
	//创建一个买单
	Convey("新增一个限价买单", t, func() {
		btcusdt.Clean()
		btcusdt.AddItem(matching.NewBidLimitItem("id11", d(1.1), d(1.2), 1112))

		So(btcusdt.AskQueue().Len(), ShouldEqual, 0)
		So(btcusdt.BidQueue().Len(), ShouldEqual, 1)
		So(btcusdt.BidQueue().Top().GetPrice(), ShouldEqual, d(1.1))
		So(btcusdt.BidQueue().Top().GetUniqueId(), ShouldEqual, "id11")
		So(btcusdt.BidQueue().Top().GetQuantity(), ShouldEqual, d(1.2))
		So(btcusdt.BidQueue().Top().GetCreateTime(), ShouldEqual, 1112)
	})

	Convey("新增一个限价卖单", t, func() {
		btcusdt.Clean()
		btcusdt.AddItem(matching.NewAskLimitItem("id12", d(1.1), d(1.2), 1112))

		So(btcusdt.AskQueue().Len(), ShouldEqual, 1)
		So(btcusdt.BidQueue().Len(), ShouldEqual, 0)
		So(btcusdt.AskQueue().Top().GetPrice(), ShouldEqual, d(1.1))
		So(btcusdt.AskQueue().Top().GetUniqueId(), ShouldEqual, "id12")
		So(btcusdt.AskQueue().Top().GetQuantity(), ShouldEqual, d(1.2))
		So(btcusdt.AskQueue().Top().GetCreateTime(), ShouldEqual, 1112)

	})

	Convey("清空买卖队列", t, func() {
		btcusdt.Clean()
		So(btcusdt.AskQueue().Len(), ShouldEqual, 0)
		So(btcusdt.BidQueue().Len(), ShouldEqual, 0)
	})

	Convey("限价买卖单，价格一致，完全成交", t, func() {
		btcusdt.Clean()

		var tradeLogs []types.TradeResult
		btcusdt.OnTradeResult(func(result types.TradeResult) {
			tradeLogs = append(tradeLogs, result)
		})

		btcusdt.AddItem(matching.NewAskLimitItem("id13", d(1.1), d(1.2), time.Now().UnixNano()))
		btcusdt.AddItem(matching.NewBidLimitItem("id23", d(1.1), d(1.2), time.Now().UnixNano()))

		time.Sleep(time.Second)
		So(tradeLogs[0].AskOrderId, ShouldEqual, "id13")
		So(tradeLogs[0].BidOrderId, ShouldEqual, "id23")
		So(tradeLogs[0].TradePrice, ShouldEqual, d(1.1))
		So(tradeLogs[0].TradeQuantity, ShouldEqual, d(1.2))
		//ask提供流动性，bid主动成交 trade_by=buyer
		So(tradeLogs[0].TradeBy, ShouldEqual, types.ByBuyer)
	})

	Convey("限价买卖单，价格一致，买单部分成交", t, func() {
		btcusdt.Clean()

		var tradeLogs []types.TradeResult
		btcusdt.OnTradeResult(func(result types.TradeResult) {
			tradeLogs = append(tradeLogs, result)
		})

		btcusdt.AddItem(matching.NewBidLimitItem("id24", d(1.1), d(2.3), time.Now().UnixNano()))
		btcusdt.AddItem(matching.NewAskLimitItem("id14", d(1.1), d(1.2), time.Now().UnixNano()))
		time.Sleep(time.Second)

		So(tradeLogs[0].AskOrderId, ShouldEqual, "id14")
		So(tradeLogs[0].BidOrderId, ShouldEqual, "id24")
		So(tradeLogs[0].TradePrice, ShouldEqual, d(1.1))
		So(tradeLogs[0].TradeQuantity, ShouldEqual, d(1.2))
		So(btcusdt.BidQueue().Top().GetQuantity(), ShouldEqual, d(1.1))
		So(btcusdt.BidQueue().Len(), ShouldEqual, 1)
		So(btcusdt.AskQueue().Len(), ShouldEqual, 0)
		//bid提供流动性，ask主动成交 trade_by=seller
		So(tradeLogs[0].TradeBy, ShouldEqual, types.BySeller)
	})

	Convey("限价买卖单，价格一致，卖单部分成交", t, func() {
		btcusdt.Clean()

		var tradeLogs []types.TradeResult
		btcusdt.OnTradeResult(func(result types.TradeResult) {
			tradeLogs = append(tradeLogs, result)
		})

		btcusdt.AddItem(matching.NewAskLimitItem("id15", d(1.1), d(2.2), 1112))
		btcusdt.AddItem(matching.NewBidLimitItem("id25", d(1.1), d(1.3), 1113))
		time.Sleep(time.Second)

		So(tradeLogs[0].AskOrderId, ShouldEqual, "id15")
		So(tradeLogs[0].BidOrderId, ShouldEqual, "id25")
		So(tradeLogs[0].TradePrice, ShouldEqual, d(1.1))
		So(tradeLogs[0].TradeQuantity, ShouldEqual, d(1.3))
		So(btcusdt.AskQueue().Top().GetQuantity(), ShouldEqual, d(0.9))
		So(btcusdt.AskQueue().Len(), ShouldEqual, 1)
		So(btcusdt.BidQueue().Len(), ShouldEqual, 0)
	})

	Convey("价格一致，时间优先", t, func() {
		btcusdt.Clean()

		var tradeLogs []types.TradeResult
		btcusdt.OnTradeResult(func(result types.TradeResult) {
			tradeLogs = append(tradeLogs, result)
		})

		btcusdt.AddItem(matching.NewAskLimitItem("id16", d(1.1), d(2.2), 1112))
		btcusdt.AddItem(matching.NewAskLimitItem("id26", d(1.1), d(2.2), 1110))

		btcusdt.AddItem(matching.NewBidLimitItem("id36", d(1.1), d(1.3), 1113))
		time.Sleep(time.Second)

		So(tradeLogs[0].AskOrderId, ShouldEqual, "id26")
		So(tradeLogs[0].BidOrderId, ShouldEqual, "id36")
		So(tradeLogs[0].TradePrice, ShouldEqual, d(1.1))
		So(tradeLogs[0].TradeQuantity, ShouldEqual, d(1.3))
		So(btcusdt.AskQueue().Top().GetQuantity(), ShouldEqual, d(0.9))
		So(btcusdt.AskQueue().Len(), ShouldEqual, 2)
		So(btcusdt.BidQueue().Len(), ShouldEqual, 0)
	})

	//价格优先
	Convey("价格优先", t, func() {
		btcusdt.Clean()

		var tradeLogs []types.TradeResult
		btcusdt.OnTradeResult(func(result types.TradeResult) {
			tradeLogs = append(tradeLogs, result)
		})

		btcusdt.AddItem(matching.NewAskLimitItem("id17", d(1.01), d(2.2), 1112))
		btcusdt.AddItem(matching.NewAskLimitItem("id27", d(1.1), d(2.2), 1110))

		btcusdt.AddItem(matching.NewBidLimitItem("id37", d(1.1), d(1.3), 1113))
		time.Sleep(time.Second)

		So(tradeLogs[0].AskOrderId, ShouldEqual, "id17")
		So(tradeLogs[0].BidOrderId, ShouldEqual, "id37")
		So(tradeLogs[0].TradePrice, ShouldEqual, d(1.01))
		So(tradeLogs[0].TradeQuantity, ShouldEqual, d(1.3))
		So(tradeLogs[0].RemainderMarketOrderId, ShouldEqual, "")
		So(btcusdt.AskQueue().Top().GetQuantity(), ShouldEqual, d(0.9))
		So(btcusdt.AskQueue().Len(), ShouldEqual, 2)
		So(btcusdt.BidQueue().Len(), ShouldEqual, 0)
	})

}

func TestTradeFunc_MarketBuyOrder(t *testing.T) {

	Convey("市价买入 按数量买入, 金额足够买单完全成交", t, func() {
		btcusdt.Clean()

		var tradeLogs []types.TradeResult
		btcusdt.OnTradeResult(func(result types.TradeResult) {
			tradeLogs = append(tradeLogs, result)
		})

		btcusdt.AddItem(matching.NewAskLimitItem("id18", d(1.01), d(2.2), 1112))
		btcusdt.AddItem(matching.NewBidMarketQtyItem("id28", d(1.1), d(100), 1113))
		time.Sleep(time.Second)

		So(tradeLogs[0].AskOrderId, ShouldEqual, "id18")
		So(tradeLogs[0].BidOrderId, ShouldEqual, "id28")
		So(tradeLogs[0].TradePrice, ShouldEqual, d(1.01))
		So(tradeLogs[0].TradeQuantity, ShouldEqual, d(1.1))
		So(tradeLogs[0].RemainderMarketOrderId, ShouldEqual, "id28")

	})

	Convey("市价按数量买入, 账户金额足够买单部分成交", t, func() {
		btcusdt.Clean()

		var tradeLogs []types.TradeResult
		btcusdt.OnTradeResult(func(result types.TradeResult) {
			tradeLogs = append(tradeLogs, result)
		})

		btcusdt.AddItem(matching.NewAskLimitItem("id19", d(1.01), d(2.2), 1112))
		btcusdt.AddItem(matching.NewBidMarketQtyItem("id29", d(100), d(100), 1113))
		time.Sleep(time.Second)

		So(tradeLogs[0].AskOrderId, ShouldEqual, "id19")
		So(tradeLogs[0].BidOrderId, ShouldEqual, "id29")
		So(tradeLogs[0].TradePrice, ShouldEqual, d(1.01))
		So(tradeLogs[0].TradeQuantity, ShouldEqual, d(2.2))
		So(tradeLogs[0].RemainderMarketOrderId, ShouldEqual, "id29")
	})

	//市价买入 按数量, 金额不足 买单部分成交
	Convey("市价按数量买入, 账户金额不足, 买单部分成交", t, func() {
		btcusdt.Clean()

		var tradeLogs []types.TradeResult
		btcusdt.OnTradeResult(func(result types.TradeResult) {
			tradeLogs = append(tradeLogs, result)
		})

		btcusdt.AddItem(matching.NewAskLimitItem("id110", d(100), d(20), 1112))
		btcusdt.AddItem(matching.NewBidMarketQtyItem("id210", d(20), d(100), 1113))
		time.Sleep(time.Second)

		So(tradeLogs[0].AskOrderId, ShouldEqual, "id110")
		So(tradeLogs[0].BidOrderId, ShouldEqual, "id210")
		So(tradeLogs[0].TradePrice, ShouldEqual, d(100))
		So(tradeLogs[0].TradeQuantity, ShouldEqual, d(1))

		So(tradeLogs[0].RemainderMarketOrderId, ShouldEqual, "id210")
		So(btcusdt.AskQueue().Len(), ShouldEqual, 1)
		So(btcusdt.BidQueue().Len(), ShouldEqual, 0)
	})

	Convey("市价买入 指定金额, 买单完全成交", t, func() {
		btcusdt.Clean()

		var tradeLogs []types.TradeResult
		btcusdt.OnTradeResult(func(result types.TradeResult) {
			tradeLogs = append(tradeLogs, result)
		})

		btcusdt.AddItem(matching.NewAskLimitItem("id111", d(10.00), d(100), 1112))
		btcusdt.AddItem(matching.NewBidMarketAmountItem("id211", d(50), 1113))
		time.Sleep(time.Second)

		So(tradeLogs[0].AskOrderId, ShouldEqual, "id111")
		So(tradeLogs[0].BidOrderId, ShouldEqual, "id211")
		So(tradeLogs[0].TradePrice, ShouldEqual, d(10.00))
		So(tradeLogs[0].TradeQuantity, ShouldEqual, d(5))

		So(tradeLogs[0].RemainderMarketOrderId, ShouldEqual, "id211")
		So(btcusdt.AskQueue().Len(), ShouldEqual, 1)
	})

	Convey("市价买入 指定金额, 买单部分成交", t, func() {
		btcusdt.Clean()

		var tradeLogs []types.TradeResult
		btcusdt.OnTradeResult(func(result types.TradeResult) {
			tradeLogs = append(tradeLogs, result)
		})

		btcusdt.AddItem(matching.NewAskLimitItem("id112", d(10.00), d(100), 1112))
		btcusdt.AddItem(matching.NewBidMarketAmountItem("id212", d(6000), 1113))
		time.Sleep(time.Second)

		So(tradeLogs[0].AskOrderId, ShouldEqual, "id112")
		So(tradeLogs[0].BidOrderId, ShouldEqual, "id212")
		So(tradeLogs[0].TradePrice, ShouldEqual, d(10.00))
		So(tradeLogs[0].TradeQuantity, ShouldEqual, d(100))

		So(tradeLogs[0].RemainderMarketOrderId, ShouldEqual, "id212")
		So(btcusdt.AskQueue().Len(), ShouldEqual, 0)
		So(btcusdt.BidQueue().Len(), ShouldEqual, 0)
	})
}

func TestTradeFunc_MarketSellOrder(t *testing.T) {
	Convey("市价卖出 按数量, 完全成交", t, func() {

		btcusdt.Clean()

		var tradeLogs []types.TradeResult
		btcusdt.OnTradeResult(func(result types.TradeResult) {
			tradeLogs = append(tradeLogs, result)
		})

		btcusdt.AddItem(matching.NewBidLimitItem("id113", d(10.00), d(100), 1112))
		btcusdt.AddItem(matching.NewAskMarketQtyItem("id213", d(6), 1113))
		time.Sleep(time.Second)

		So(tradeLogs[0].AskOrderId, ShouldEqual, "id213")
		So(tradeLogs[0].BidOrderId, ShouldEqual, "id113")
		So(tradeLogs[0].TradePrice, ShouldEqual, d(10.00))
		So(tradeLogs[0].TradeQuantity, ShouldEqual, d(6))

		So(tradeLogs[0].RemainderMarketOrderId, ShouldEqual, "id213")
		So(btcusdt.AskQueue().Len(), ShouldEqual, 0)
		So(btcusdt.BidQueue().Len(), ShouldEqual, 1)
	})

	Convey("市价卖出 按数量，部分成交", t, func() {
		//
		btcusdt.Clean()

		var tradeLogs []types.TradeResult
		btcusdt.OnTradeResult(func(result types.TradeResult) {
			tradeLogs = append(tradeLogs, result)
		})

		btcusdt.AddItem(matching.NewBidLimitItem("id114", d(10.00), d(100), 1112))
		btcusdt.AddItem(matching.NewAskMarketQtyItem("id214", d(6000), 1113))
		time.Sleep(time.Second)

		So(tradeLogs[0].AskOrderId, ShouldEqual, "id214")
		So(tradeLogs[0].BidOrderId, ShouldEqual, "id114")
		So(tradeLogs[0].TradePrice, ShouldEqual, d(10.00))
		So(tradeLogs[0].TradeQuantity, ShouldEqual, d(100))

		So(tradeLogs[0].RemainderMarketOrderId, ShouldEqual, "id214")
		So(btcusdt.AskQueue().Len(), ShouldEqual, 0)
		So(btcusdt.BidQueue().Len(), ShouldEqual, 0)
	})

	Convey("市价卖出 指定金额，持仓足够完全成交", t, func() {
		btcusdt.Clean()

		var tradeLogs []types.TradeResult
		btcusdt.OnTradeResult(func(result types.TradeResult) {
			tradeLogs = append(tradeLogs, result)
		})

		btcusdt.AddItem(matching.NewBidLimitItem("id115", d(10.00), d(1000), 1112))
		btcusdt.AddItem(matching.NewAskMarketAmountItem("id215", d(6000), d(1000000), 1113))
		time.Sleep(time.Second)

		So(tradeLogs[0].AskOrderId, ShouldEqual, "id215")
		So(tradeLogs[0].BidOrderId, ShouldEqual, "id115")
		So(tradeLogs[0].TradePrice, ShouldEqual, d(10.00))
		So(tradeLogs[0].TradeQuantity, ShouldEqual, d(600))

		So(tradeLogs[0].RemainderMarketOrderId, ShouldEqual, "id215")
		So(btcusdt.AskQueue().Len(), ShouldEqual, 0)
		So(btcusdt.BidQueue().Len(), ShouldEqual, 1)
	})

	Convey("市价卖出 指定金额，持仓足够 但是部分成交", t, func() {
		btcusdt.Clean()

		var tradeLogs []types.TradeResult
		btcusdt.OnTradeResult(func(result types.TradeResult) {
			tradeLogs = append(tradeLogs, result)
		})

		btcusdt.AddItem(matching.NewBidLimitItem("id116", d(10.00), d(50), 1112))
		btcusdt.AddItem(matching.NewAskMarketAmountItem("id216", d(6000), d(1000000), 1113))
		time.Sleep(time.Second)

		So(tradeLogs[0].AskOrderId, ShouldEqual, "id216")
		So(tradeLogs[0].BidOrderId, ShouldEqual, "id116")
		So(tradeLogs[0].TradePrice, ShouldEqual, d(10.00))
		So(tradeLogs[0].TradeQuantity, ShouldEqual, d(50))

		So(tradeLogs[0].RemainderMarketOrderId, ShouldEqual, "id216")
		So(btcusdt.AskQueue().Len(), ShouldEqual, 0)
		So(btcusdt.BidQueue().Len(), ShouldEqual, 0)
	})

	Convey("市价卖出 指定金额，持仓不足", t, func() {
		btcusdt.Clean()

		var tradeLogs []types.TradeResult
		btcusdt.OnTradeResult(func(result types.TradeResult) {
			tradeLogs = append(tradeLogs, result)
		})

		btcusdt.AddItem(matching.NewBidLimitItem("id117", d(100.00), d(50), 1112))
		btcusdt.AddItem(matching.NewAskMarketAmountItem("id217", d(500), d(3), 1113))
		time.Sleep(time.Second)

		So(tradeLogs[0].AskOrderId, ShouldEqual, "id217")
		So(tradeLogs[0].BidOrderId, ShouldEqual, "id117")
		So(tradeLogs[0].TradePrice, ShouldEqual, d(100))
		So(tradeLogs[0].TradeQuantity, ShouldEqual, d(3))

		So(tradeLogs[0].RemainderMarketOrderId, ShouldEqual, "id217")
		So(btcusdt.AskQueue().Len(), ShouldEqual, 0)
		So(btcusdt.BidQueue().Len(), ShouldEqual, 1)
	})

	Convey("市价卖出 指定金额，一个都没有成交，只输出市价撤单信号", t, func() {
		// btcusdt.Clean()

		// cancel := make(chan matching.CancelBody)
		// go func() {
		// 	for {
		// 		select {
		// 		case cancel <- <-btcusdt.ChCancelResult:

		// 		}
		// 	}
		// }()

		// btcusdt.AddItem(matching.NewBidLimitItem("id1118", d(1000.00), d(50), 1112))
		// btcusdt.AddItem(matching.NewAskMarketAmountItem("id2218", d(1), d(30), 1113))

		// So(btcusdt.BidQueue().Len(), ShouldEqual, 1)
		// bid := btcusdt.BidQueue().Top()
		// So(bid.GetQuantity(), ShouldEqual, d(50))
		// So(btcusdt.AskQueue().Len(), ShouldEqual, 0)

		// cinfo := <-cancel
		// So(cinfo.OrderId, ShouldEqual, "id2218")
	})
}
