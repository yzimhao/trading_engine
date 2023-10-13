package assets

// import (
// 	_ "github.com/go-sql-driver/mysql"
// 	_ "github.com/lib/pq"

// 	"github.com/sirupsen/logrus"
// 	. "github.com/smartystreets/goconvey/convey"
// 	"xorm.io/xorm"
// )

// func init() {
// 	driver := "mysql"
// 	dsn := "root:root@tcp(localhost:13306)/test?charset=utf8&loc=Local"

// 	logrus.Infof("dsn: %s", dsn)

// 	conn, err := xorm.NewEngine(driver, dsn)
// 	if err != nil {
// 		logrus.Panic(err)
// 	}
// 	db_engine = conn
// 	db_engine.ShowSQL(true)

// 	db_engine.DropTables(
// 		new(Assets),
// 		new(assetsLog),
// 		new(assetFreezeRecord),
// 	)

// 	Init(db_engine, nil)
// }

// func cleanUserAssets(user_id int64) (err error) {
// 	db := db_engine.NewSession()
// 	defer db.Close()

// 	db.Begin()
// 	defer func() {
// 		if err != nil {
// 			db.Rollback()
// 		} else {
// 			db.Commit()
// 		}
// 	}()

// 	//资产表
// 	_, err = db.Table(new(Assets)).Where("user_id=?", user_id).Delete()
// 	if err != nil {
// 		return err
// 	}

// 	//资产变化日志表
// 	_, err = db.Table(new(assetsLog)).Where("user_id=?", user_id).Delete()
// 	if err != nil {
// 		return err
// 	}
// 	//资产冻结记录
// 	_, err = db.Table(new(assetFreezeRecord)).Where("user_id=?", user_id).Delete()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func Test_main(t *testing.T) {
// 	db := db_engine.NewSession()
// 	defer db.Close()

// 	var (
// 		user1 int64 = 101
// 		user2 int64 = 102

// 		symbol_usd int = 1
// 		symbol_eth int = 2
// 	)

// 	initAssets := func(uid int64, sid int, amount string, buid string) {
// 		InitAssetsForDemo(uid, sid, amount, buid)
// 	}

// 	Convey("从根账户充值100", t, func() {
// 		initAssets(user1, symbol_usd, "100", "T01")
// 		user := UserAssets(user1, symbol_usd)
// 		So(d(user.Total), ShouldEqual, d("100"))
// 		So(d(user.Available), ShouldEqual, d("100"))
// 		So(d(user.Freeze), ShouldEqual, d("0"))

// 		err := cleanUserAssets(user1)
// 		So(err, ShouldBeNil)
// 	})

// 	Convey("自己给自己转账", t, func() {
// 		before := UserAssets(user1, symbol_usd)
// 		_, err := transfer(db, true, user1, user1, symbol_usd, "10.00", "t_10", Behavior_Transfer)
// 		So(err, ShouldBeError, fmt.Errorf("invalid to"))

// 		after := UserAssets(user1, symbol_usd)
// 		So(before.Total, ShouldEqual, after.Total)

// 	})

// 	Convey("冻结用户资产", t, func() {
// 		initAssets(user1, symbol_usd, "100", "T02")
// 		f, err := freezeAssets(db, true, user1, symbol_usd, "10", "a001", Behavior_Trade)
// 		So(err, ShouldBeNil)
// 		So(f, ShouldBeTrue)

// 		user := UserAssets(user1, symbol_usd)
// 		So(d(user.Total), ShouldEqual, d("100"))
// 		So(d(user.Available), ShouldEqual, d("90"))
// 		So(d(user.Freeze), ShouldEqual, d("10"))

// 		err = cleanUserAssets(user1)
// 		So(err, ShouldBeNil)
// 	})

// 	Convey("冻结负数的资产", t, func() {
// 		initAssets(user1, symbol_usd, "100", "T03")

// 		f, err := freezeAssets(db, true, user1, symbol_usd, "-10", "a002", Behavior_Trade)
// 		So(err, ShouldBeError, fmt.Errorf("freeze amount should be >= 0"))
// 		So(f, ShouldBeFalse)

// 		err = cleanUserAssets(user1)
// 		So(err, ShouldBeNil)
// 	})

// 	Convey("冻结数量0的资产，则冻结全部可用", t, func() {
// 		initAssets(user1, symbol_usd, "100", "T04")
// 		f, err := freezeAssets(db, true, user1, symbol_usd, "0", "a003", Behavior_Trade)
// 		So(err, ShouldBeNil)
// 		So(f, ShouldBeTrue)

// 		user := UserAssets(user1, symbol_usd)
// 		So(d(user.Total), ShouldEqual, d("100"))
// 		So(d(user.Available), ShouldEqual, d("0"))
// 		So(d(user.Freeze), ShouldEqual, d("100"))

// 		err = cleanUserAssets(user1)
// 		So(err, ShouldBeNil)
// 	})

// 	Convey("解冻不存在的业务订单号", t, func() {
// 		f, err := unfreezeAssets(db, true, user1, "a004", "10")
// 		So(err, ShouldBeError, fmt.Errorf("not found business_id"))
// 		So(f, ShouldBeFalse)
// 	})

// 	Convey("解冻订单号剩余全部资产", t, func() {
// 		initAssets(user1, symbol_usd, "100", "T05")
// 		freezeAssets(db, true, user1, symbol_usd, "10", "a005", Behavior_Trade)
// 		f, err := unfreezeAssets(db, true, user1, "a005", "0")
// 		So(err, ShouldBeNil)
// 		So(f, ShouldBeTrue)

// 		user := UserAssets(user1, symbol_usd)
// 		So(d(user.Total), ShouldEqual, d("100"))
// 		So(d(user.Available), ShouldEqual, d("100"))
// 		So(d(user.Freeze), ShouldEqual, d("0"))

// 		err = cleanUserAssets(user1)
// 		So(err, ShouldBeNil)
// 	})

// 	Convey("解冻业务订单部分资产", t, func() {
// 		initAssets(user1, symbol_usd, "100", "T06")
// 		freezeAssets(db, true, user1, symbol_usd, "10", "a006", Behavior_Trade)
// 		f, err := unfreezeAssets(db, true, user1, "a006", "1.2")
// 		So(err, ShouldBeNil)
// 		So(f, ShouldBeTrue)

// 		user := UserAssets(user1, symbol_usd)
// 		So(d(user.Total), ShouldEqual, d("100"))
// 		So(d(user.Available), ShouldEqual, d("91.2"))
// 		So(d(user.Freeze), ShouldEqual, d("8.8"))

// 		err = cleanUserAssets(user1)
// 		So(err, ShouldBeNil)
// 	})

// 	Convey("解冻超过业务订单金额的数量", t, func() {
// 		initAssets(user1, symbol_usd, "100", "T07")
// 		freezeAssets(db, true, user1, symbol_usd, "10", "a006", Behavior_Trade)
// 		f, err := unfreezeAssets(db, true, user1, "a006", "11")
// 		So(err, ShouldBeError, fmt.Errorf("unfreeze amount must lt freeze amount"))
// 		So(f, ShouldBeFalse)

// 		err = cleanUserAssets(user1)
// 		So(err, ShouldBeNil)
// 	})

// 	Convey("解冻负数的数量", t, func() {
// 		initAssets(user1, symbol_usd, "100", "T08")
// 		freezeAssets(db, true, user1, symbol_usd, "10", "a006", Behavior_Trade)
// 		f, err := unfreezeAssets(db, true, user1, "a006", "-2")
// 		So(err, ShouldBeError, fmt.Errorf("unfreeze amount should be >= 0"))
// 		So(f, ShouldBeFalse)

// 		err = cleanUserAssets(user1)
// 		So(err, ShouldBeNil)
// 	})

// 	//上面测试做完后，生成一点测试数据
// 	cleanUserAssets(ROOTUSERID)
// 	initAssets(user1, symbol_usd, "1000000", "T09")
// 	initAssets(user1, symbol_eth, "1000000", "T10")

// 	initAssets(user2, symbol_usd, "1000000", "T11")
// 	initAssets(user2, symbol_eth, "1000000", "T12")
// }
