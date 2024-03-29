package assets

import (
	"github.com/yzimhao/trading_engine/utils/app"
)

const (
	//根账户，所有的进出账都从这个账户开始，所有列的资产求和应该为0
	UserRoot string = "root"
	//系统收取手续费用放在该用户ID下
	UserSystemFee string = "fee"
)

func UserAssets(user_id string, symbol []string) []Assets {
	db_engine := app.Database().NewSession()
	defer db_engine.Close()

	rows := []Assets{}
	q := db_engine.Table(new(Assets)).Where("user_id=?", user_id)
	if len(symbol) > 0 {
		q = q.In("symbol", symbol)
	}

	q.Find(&rows)
	return rows
}

func DemoData() {
	SysDeposit("user1", "usd", "10000.00", "C001")
	SysDeposit("user1", "jpy", "10000.00", "C002")
	SysDeposit("user2", "usd", "10000.00", "C001")
	SysDeposit("user2", "jpy", "10000.00", "C002")
}
