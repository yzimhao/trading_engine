package assets

import (
	"fmt"
	"strings"

	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
	"xorm.io/xorm"
)

func Transfer(db *xorm.Session, from, to string, symbol string, amount string, business_id string, behavior OpBehavior) (success bool, err error) {
	return transfer(db, from, to, symbol, amount, business_id, behavior)
}

// 充值
func SysDeposit(to string, symbol string, amount string, business_id string) (success bool, err error) {
	db := app.Database().NewSession()
	defer db.Close()

	db.Begin()
	success, err = transfer(db, UserRoot, to, symbol, amount, business_id, Behavior_Recharge)
	if err != nil {
		db.Rollback()
	}
	db.Commit()
	return success, err
}

// 提现
func SysWithdraw(user_id string, symbol string, amount string, business_id string) (success bool, err error) {
	db := app.Database().NewSession()
	defer db.Close()

	db.Begin()
	success, err = transfer(db, user_id, UserRoot, symbol, amount, business_id, Behavior_Withdraw)
	if err != nil {
		db.Rollback()
	}
	db.Commit()
	return success, err
}

func transfer(db *xorm.Session, from, to string, symbol string, amount string, business_id string, behavior OpBehavior) (success bool, err error) {
	symbol = strings.ToLower(symbol)

	if from == to {
		return false, fmt.Errorf("invalid to")
	}

	from_user := Assets{UserId: from, Symbol: symbol}
	has_from, err := db.Table(new(Assets)).ForUpdate().Get(&from_user)
	if err != nil {
		return false, err
	}

	//非根账户检查余额
	if from != UserRoot {
		if utils.D(from_user.Available).Cmp(utils.D("0")) <= 0 {
			return false, fmt.Errorf("可用资产不足")
		}
	}

	to_user := Assets{UserId: to, Symbol: symbol}
	has_to, err := db.Table(new(Assets)).ForUpdate().Get(&to_user)
	if err != nil {
		return false, err
	}

	from_before := utils.D(from_user.Total)
	from_user.Total = utils.D(from_user.Total).Sub(utils.D(amount)).String()
	from_user.Available = utils.D(from_user.Available).Sub(utils.D(amount)).String()
	if !has_from {
		from_user.Freeze = "0"
		_, err = db.Table(new(Assets)).Insert(&from_user)
	} else {
		_, err = db.Table(new(Assets)).Where("user_id=? and symbol=?", from, symbol).Update(&from_user)
	}
	if err != nil {
		return false, err
	}

	to_before := utils.D(to_user.Total)
	to_user.Total = utils.D(to_user.Total).Add(utils.D(amount)).String()
	to_user.Available = utils.D(to_user.Available).Add(utils.D(amount)).String()
	if !has_to {
		to_user.Freeze = "0"
		_, err = db.Table(new(Assets)).Insert(&to_user)
	} else {
		_, err = db.Table(new(Assets)).Where("user_id=? and symbol=?", to, symbol).Update(&to_user)
	}
	if err != nil {
		return false, err
	}

	//双方日志
	from_log := assetsLog{
		UserId:     from,
		Symbol:     symbol,
		Before:     from_before.String(),
		Amount:     "-" + amount,
		After:      from_user.Total,
		BusinessId: business_id,
		Behavior:   behavior,
		Info:       fmt.Sprintf("to: %s", to),
	}
	_, err = db.Table(new(assetsLog)).Insert(&from_log)
	if err != nil {
		return false, err
	}

	to_log := assetsLog{
		UserId:     to,
		Symbol:     symbol,
		Before:     to_before.String(),
		Amount:     amount,
		After:      to_user.Total,
		BusinessId: business_id,
		Behavior:   behavior,
		Info:       fmt.Sprintf("from: %s", from),
	}
	_, err = db.Table(new(assetsLog)).Insert(&to_log)
	if err != nil {
		return false, err
	}
	return true, err
}
