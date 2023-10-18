package orders

import (
	"errors"
	"testing"

	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haobase/base/symbols"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/utils/app"

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

func initdb(t *testing.T) {
	app.DatabaseInit("mysql", "root:root@tcp(localhost:3306)/test?charset=utf8&loc=Local", true)
	app.RedisInit("127.0.0.1:6379", "", 15)

	cleanSymbols(t)
	cleanAssets(t)
	cleanOrders(t)
	base.Init()
}

func initAssets(t *testing.T) {
	assets.Init()
	symbols.DemoData()

	assets.SysRecharge(sellUser, testTargetSymbol, "10000.00", "C001")
	assets.SysRecharge(buyUser, testBaseSymbol, "10000.00", "C001")
}

func cleanAssets(t *testing.T) {
	db := app.Database()
	db.DropIndexes(new(assets.Assets))
	db.DropIndexes("assets_freeze")
	db.DropIndexes("assets_log")
	err := db.DropTables(new(assets.Assets), "assets_freeze", "assets_log")
	if err != nil {
		t.Logf("mysql droptables: %s", err)
	}

}

func cleanSymbols(t *testing.T) {
	db := app.Database()
	db.DropIndexes(new(symbols.Varieties))
	db.DropIndexes(new(symbols.TradingVarieties))
	err := db.DropTables(new(symbols.Varieties), new(symbols.TradingVarieties))
	if err != nil {
		t.Logf("mysql droptables: %s", err)
	}
}

func cleanOrders(t *testing.T) {
	db := app.Database()

	db.DropIndexes(new(UnfinishedOrder))
	db.DropIndexes(GetOrderTableName(testSymbol))
	err := db.DropTables(new(UnfinishedOrder), GetOrderTableName(testSymbol))
	if err != nil {
		t.Logf("mysql droptables: %s", err)
	}
}

func TestNewOrder(t *testing.T) {
	initdb(t)
	initAssets(t)

	defer cleanOrders(t)
	defer cleanAssets(t)

	Convey("新订单下单测试", t, func() {
		_, err := NewLimitOrder(sellUser, testSymbol, trading_core.OrderSideSell, "1.00", "1")
		So(err, ShouldBeNil)

		_, err = NewLimitOrder(buyUser, testSymbol, trading_core.OrderSideBuy, "1.00", "1")
		So(err, ShouldBeNil)
	})
}

func TestNewOrderCase1(t *testing.T) {
	initdb(t)
	initAssets(t)

	// defer cleanOrders(t)
	// defer cleanAssets(t)

	Convey("用户反向有挂单 测试新开订单", t, func() {
		assets.SysRecharge(sellUser, testBaseSymbol, "10000.00", "C001")
		_, err := limit_order(sellUser, testSymbol, trading_core.OrderSideSell, "1.00", "1")
		So(err, ShouldBeNil)

		_, err = limit_order(sellUser, testSymbol, trading_core.OrderSideBuy, "1.00", "1")
		So(err, ShouldBeError, errors.New("对向有挂单请撤单后再操作"))

		_, err = limit_order(sellUser, testSymbol, trading_core.OrderSideBuy, "0.9", "1")
		So(err, ShouldBeNil)
	})
}
