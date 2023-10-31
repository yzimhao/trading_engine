package clearing

import (
	"fmt"
	"testing"
	"time"

	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/cmd/haobase/base/varieties"
	"github.com/yzimhao/trading_engine/cmd/haobase/orders"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
	"xorm.io/xorm/log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	sellUser         = "seller1"
	buyUser          = "buyer1"
	testSymbol       = "usdjpy"
	testTargetSymbol = "usd"
	testBaseSymbol   = "jpy"
)

func init() {
	initdb()
}

func initdb() {
	app.ConfigInit("", false)
	app.DatabaseInit("mysql", "root:root@tcp(localhost:3306)/test?charset=utf8&loc=Local", true, "")
	app.Database().SetLogLevel(log.LOG_DEBUG)
	app.RedisInit("127.0.0.1:6379", "", 15)
	cleanSymbols()
	initSymbols()
}

func initSymbols() {
	varieties.Init()
	varieties.DemoData()
}

func cleandb() {
	cleanAssets()
	cleanOrders()
}

func initAssets() {
	assets.Init()
	assets.SysDeposit(sellUser, testTargetSymbol, "10000.00", "C001")
	assets.SysDeposit(buyUser, testBaseSymbol, "10000.00", "C001")
}

func cleanAssets() {
	db := app.Database()
	db.DropIndexes(new(assets.Assets))
	db.DropIndexes("assets_freeze")
	db.DropIndexes("assets_log")
	err := db.DropTables(new(assets.Assets), "assets_freeze", "assets_log")
	if err != nil {
		fmt.Errorf("mysql droptables: %s", err)
	}

}

func cleanSymbols() {
	db := app.Database()
	db.DropIndexes(new(varieties.Varieties))
	db.DropIndexes(new(varieties.TradingVarieties))
	err := db.DropTables(new(varieties.Varieties), new(varieties.TradingVarieties))
	if err != nil {
		fmt.Errorf("mysql droptables: %s", err)
	}
}

func cleanOrders() {
	db := app.Database()
	db.DropIndexes(orders.GetOrderTableName(testSymbol))
	db.DropIndexes(new(orders.UnfinishedOrder))
	db.DropIndexes(orders.GetTradelogTableName(testSymbol))

	db.DropTables(orders.GetOrderTableName(testSymbol))
	db.DropTables(new(orders.UnfinishedOrder))
	db.DropTables(orders.GetTradelogTableName(testSymbol))

}

func TestLimitOrder(t *testing.T) {
	cleandb()
	initAssets()

	Convey("限价单完全成交结算测试", t, func() {

		sell, err := orders.NewLimitOrder(sellUser, testSymbol, trading_core.OrderSideSell, "1.00", "1")
		So(err, ShouldBeNil)

		buy, err := orders.NewLimitOrder(buyUser, testSymbol, trading_core.OrderSideBuy, "1.00", "1")
		So(err, ShouldBeNil)

		result := trading_core.TradeResult{
			Symbol:        testSymbol,
			AskOrderId:    sell.OrderId,
			BidOrderId:    buy.OrderId,
			TradePrice:    utils.D("1.00"),
			TradeQuantity: utils.D("1"),
			TradeTime:     time.Now().UnixNano(),
		}
		clearing_trade_order(testSymbol, result.Json())

		for {
			if orders.GetLock(orders.ClearingLock, sell.OrderId) == 0 &&
				orders.GetLock(orders.ClearingLock, buy.OrderId) == 0 {
				break
			}
			time.Sleep(1 * time.Second)
		}
		//检查资产
		sell_assets_target := assets.FindSymbol(sellUser, testTargetSymbol)
		sell_assets_standard := assets.FindSymbol(sellUser, testBaseSymbol)

		buy_assets_target := assets.FindSymbol(buyUser, testTargetSymbol)
		buy_assets_standard := assets.FindSymbol(buyUser, testBaseSymbol)
		So(utils.D(sell_assets_target.Total), ShouldEqual, utils.D("9999"))
		So(utils.D(sell_assets_standard.Total), ShouldEqual, utils.D("0.995"))

		So(utils.D(buy_assets_target.Total), ShouldEqual, utils.D("1"))
		So(utils.D(buy_assets_standard.Total), ShouldEqual, utils.D("9998.995"))

		//检查订单状态
		sell_order := orders.Find(testSymbol, sell.OrderId)
		So(sell_order.Status, ShouldEqual, orders.OrderStatusDone)
		buy_order := orders.Find(testSymbol, buy.OrderId)
		So(buy_order.Status, ShouldEqual, orders.OrderStatusDone)

		sell_unfinished := orders.FindUnfinished(testSymbol, sell.OrderId)
		So(sell_unfinished, ShouldBeNil)
		buy_unfinished := orders.FindUnfinished(testSymbol, buy.OrderId)
		So(buy_unfinished, ShouldBeNil)

	})
}

func TestMarketCase1(t *testing.T) {

	cleandb()
	initAssets()

	Convey("市价买指定的数量,完全成交", t, func() {

		s1, err := orders.NewLimitOrder(sellUser, testSymbol, trading_core.OrderSideSell, "1.00", "1")
		So(err, ShouldBeNil)

		s2, _ := orders.NewLimitOrder(sellUser, testSymbol, trading_core.OrderSideSell, "2.00", "1")

		buy, err := orders.NewMarketOrderByQty(buyUser, testSymbol, trading_core.OrderSideBuy, "3")
		So(err, ShouldBeNil)

		result1 := trading_core.TradeResult{
			Symbol:        testSymbol,
			AskOrderId:    s1.OrderId,
			BidOrderId:    buy.OrderId,
			TradePrice:    utils.D("1.00"),
			TradeQuantity: utils.D("1"),
			TradeTime:     time.Now().UnixNano(),
		}
		result2 := trading_core.TradeResult{
			Symbol:        testSymbol,
			AskOrderId:    s2.OrderId,
			BidOrderId:    buy.OrderId,
			TradePrice:    utils.D("2.00"),
			TradeQuantity: utils.D("1"),
			TradeTime:     time.Now().UnixNano(),
			Last:          buy.OrderId,
		}

		clearing_trade_order(testSymbol, result2.Json())
		clearing_trade_order(testSymbol, result1.Json())

		for {
			if orders.GetLock(orders.ClearingLock, s1.OrderId) == 0 &&
				orders.GetLock(orders.ClearingLock, s2.OrderId) == 0 &&
				orders.GetLock(orders.ClearingLock, buy.OrderId) == 0 {
				break
			}
			time.Sleep(1 * time.Second)
		}

		//检查买卖双方订单状态及资产
		s1 = orders.Find(testSymbol, s1.OrderId)
		So(s1.Status, ShouldEqual, orders.OrderStatusDone)
		s2 = orders.Find(testSymbol, s2.OrderId)
		So(s2.Status, ShouldEqual, orders.OrderStatusDone)
		buy = orders.Find(testSymbol, buy.OrderId)
		So(buy.Status, ShouldEqual, orders.OrderStatusDone)

		//资产
		sell_assets_target := assets.FindSymbol(sellUser, testTargetSymbol)
		sell_assets_standard := assets.FindSymbol(sellUser, testBaseSymbol)
		buy_assets_target := assets.FindSymbol(buyUser, testTargetSymbol)
		buy_assets_standard := assets.FindSymbol(buyUser, testBaseSymbol)
		//
		So(utils.D(sell_assets_target.Total), ShouldEqual, utils.D("9998"))
		So(utils.D(sell_assets_standard.Total), ShouldEqual, utils.D("2.985"))
		So(utils.D(sell_assets_target.Freeze), ShouldEqual, utils.D("0"))

		//市价1块买入2个usd，花费jpy
		So(utils.D(buy_assets_target.Total), ShouldEqual, utils.D("2"))
		So(utils.D(buy_assets_standard.Total), ShouldEqual, utils.D("9996.985"))
		So(utils.D(buy_assets_standard.Freeze), ShouldEqual, utils.D("0"))

		//系统收入的手续费
		fee := assets.FindSymbol(assets.UserSystemFee, testBaseSymbol)
		So(utils.D(fee.Total), ShouldEqual, utils.D(s1.Fee).Add(utils.D(s2.Fee)).Add(utils.D(buy.Fee)))

	})
}

func TestMarketCase2(t *testing.T) {
	cleandb()
	initAssets()

	Convey("市价多单测试", t, func() {

		s1, _ := orders.NewLimitOrder(sellUser, testSymbol, trading_core.OrderSideSell, "1.00", "1")
		s2, _ := orders.NewLimitOrder(sellUser, testSymbol, trading_core.OrderSideSell, "2.00", "1")
		s3, _ := orders.NewLimitOrder(sellUser, testSymbol, trading_core.OrderSideSell, "2.00", "1")
		s4, _ := orders.NewLimitOrder(sellUser, testSymbol, trading_core.OrderSideSell, "2.00", "1")

		buy, err := orders.NewMarketOrderByQty(buyUser, testSymbol, trading_core.OrderSideBuy, "5")
		So(err, ShouldBeNil)

		result1 := trading_core.TradeResult{Symbol: testSymbol, AskOrderId: s1.OrderId, BidOrderId: buy.OrderId, TradePrice: utils.D("1.00"), TradeQuantity: utils.D("1"), TradeTime: time.Now().UnixNano()}
		result2 := trading_core.TradeResult{Symbol: testSymbol, AskOrderId: s2.OrderId, BidOrderId: buy.OrderId, TradePrice: utils.D("2.00"), TradeQuantity: utils.D("1"), TradeTime: time.Now().UnixNano()}
		result3 := trading_core.TradeResult{Symbol: testSymbol, AskOrderId: s3.OrderId, BidOrderId: buy.OrderId, TradePrice: utils.D("2.00"), TradeQuantity: utils.D("1"), TradeTime: time.Now().UnixNano()}
		result4 := trading_core.TradeResult{Symbol: testSymbol, AskOrderId: s4.OrderId, BidOrderId: buy.OrderId, TradePrice: utils.D("2.00"), TradeQuantity: utils.D("1"), TradeTime: time.Now().UnixNano(), Last: buy.OrderId}

		clearing_trade_order(testSymbol, result4.Json())
		clearing_trade_order(testSymbol, result2.Json())
		clearing_trade_order(testSymbol, result1.Json())
		clearing_trade_order(testSymbol, result3.Json())

		for {
			if orders.GetLock(orders.ClearingLock, s1.OrderId) == 0 &&
				orders.GetLock(orders.ClearingLock, s2.OrderId) == 0 &&
				orders.GetLock(orders.ClearingLock, s3.OrderId) == 0 &&
				orders.GetLock(orders.ClearingLock, s4.OrderId) == 0 &&
				orders.GetLock(orders.ClearingLock, buy.OrderId) == 0 {
				break
			}
			time.Sleep(1 * time.Second)
		}

		//资产
		sell_assets_target := assets.FindSymbol(sellUser, testTargetSymbol)
		sell_assets_base := assets.FindSymbol(sellUser, testBaseSymbol)
		buy_assets_target := assets.FindSymbol(buyUser, testTargetSymbol)
		buy_assets_base := assets.FindSymbol(buyUser, testBaseSymbol)

		//卖家资产检查
		So(utils.D(sell_assets_target.Total), ShouldEqual, utils.D("9996"))
		So(utils.D(sell_assets_target.Freeze), ShouldEqual, utils.D("0"))
		So(utils.D(sell_assets_base.Total), ShouldEqual, utils.D("6.965"))
		//买家资产检查
		So(utils.D(buy_assets_target.Total), ShouldEqual, utils.D("4"))      //买5个实际市场只有4个
		So(utils.D(buy_assets_base.Total), ShouldEqual, utils.D("9992.965")) //初始本金 - （成交额 + fee）
		So(utils.D(buy_assets_base.Freeze), ShouldEqual, utils.D("0"))       //交易完成后应该全部解冻

		//系统收入的手续费
		fee := assets.FindSymbol(assets.UserSystemFee, testBaseSymbol)
		So(utils.D(fee.Total), ShouldEqual, utils.D("0.07")) //本次成交金额7，手续费费率0.005，买卖双方同时收取 7*0.005*2
	})
}
