package orders

import (
	"errors"
	"fmt"
	"testing"

	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/cmd/haobase/base/varieties"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/types/dbtables"
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
	app.DatabaseInit("mysql", "root:root@tcp(localhost:3306)/test1?charset=utf8&loc=Local", true, "")
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
	assets.SysDeposit(sellUser, testTargetSymbol, "10000.00", "C001")
	assets.SysDeposit(buyUser, testBaseSymbol, "10000.00", "C001")
}

func cleanAssets() {
	db := app.Database().NewSession()
	defer db.Close()

	dbtables.CleanTable(db, &assets.Assets{})
	dbtables.CleanTable(db, &assets.AssetsFreeze{Symbol: testBaseSymbol})
	dbtables.CleanTable(db, &assets.AssetsLog{Symbol: testBaseSymbol})
	dbtables.CleanTable(db, &assets.AssetsFreeze{Symbol: testTargetSymbol})
	dbtables.CleanTable(db, &assets.AssetsLog{Symbol: testTargetSymbol})

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

	tables := []any{
		&Order{Symbol: testSymbol},
		new(UnfinishedOrder).TableName(),
		&TradeLog{Symbol: testSymbol},
	}

	s := db.NewSession()
	defer s.Close()

	for _, table := range tables {
		dbtables.CleanTable(s, table)
	}

}

func TestNewOrder(t *testing.T) {
	cleandb()
	initAssets()

	Convey("新限价单下单", t, func() {
		_, err := NewLimitOrder(sellUser, testSymbol, trading_core.OrderSideSell, "1.00", "1")
		So(err, ShouldBeNil)

		_, err = NewLimitOrder(buyUser, testSymbol, trading_core.OrderSideBuy, "1.00", "1")
		So(err, ShouldBeNil)

		//最小交易数量限制
		_, err = NewLimitOrder(buyUser, testSymbol, trading_core.OrderSideBuy, "1.00", "0.01")
		So(err, ShouldBeNil)
		_, err = NewLimitOrder(buyUser, testSymbol, trading_core.OrderSideBuy, "1.00", "0.0001")
		So(err, ShouldBeError, errors.New("数量低于交易对最小限制"))

		//最小成交量限制

	})
}

func TestNewOrderCase1(t *testing.T) {
	cleandb()
	initAssets()

	Convey("用户反向有挂单 测试新开限价单", t, func() {
		assets.SysDeposit(sellUser, testBaseSymbol, "10000.00", "C001")

		//先挂单一个单价1.00的卖
		_, err := limit_order(sellUser, testSymbol, trading_core.OrderSideSell, "1.00", "1")
		So(err, ShouldBeNil)
		//同一用户挂一个单价1.00的限价单买
		_, err = limit_order(sellUser, testSymbol, trading_core.OrderSideBuy, "1.00", "1")
		So(err, ShouldBeError, errors.New("对向有挂单请撤单后再操作"))

		//同一用户挂一个单价2.00的限价单买
		_, err = limit_order(sellUser, testSymbol, trading_core.OrderSideBuy, "2.00", "1")
		So(err, ShouldBeError, errors.New("对向有挂单请撤单后再操作"))

		//同一用户挂一个单价0.9的限价买，这个允许挂单
		_, err = limit_order(sellUser, testSymbol, trading_core.OrderSideBuy, "0.9", "1")
		So(err, ShouldBeNil)

	})
}

func TestNewOrderCase2(t *testing.T) {
	cleandb()
	initAssets()

	Convey("用户反向有挂单 测试新开市价单", t, func() {
		assets.SysDeposit(sellUser, testBaseSymbol, "10000.00", "C001")

		//先挂单一个单价1.00的卖
		_, err := limit_order(sellUser, testSymbol, trading_core.OrderSideSell, "1.00", "1")
		So(err, ShouldBeNil)
		//同一用户挂一个市价单买
		_, err = market_order_qty(sellUser, testSymbol, trading_core.OrderSideBuy, "10")
		So(err, ShouldBeError, errors.New("对向有挂单请撤单后再操作"))
		_, err = market_order_amount(sellUser, testSymbol, trading_core.OrderSideBuy, "100.00")
		So(err, ShouldBeError, errors.New("对向有挂单请撤单后再操作"))

	})
}
