package assets

import (
	"fmt"

	"github.com/yzimhao/trading_engine/utils"
	"xorm.io/xorm"
)

func QueryFreeze(db *xorm.Session, business_id string) (*AssetsFreeze, error) {
	row := AssetsFreeze{}
	has, err := db.Table(new(AssetsFreeze)).Where("business_id=?", business_id).Get(&row)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, fmt.Errorf("冻结记录不存在")
	}

	return &row, nil
}

func FreezeAssets(db *xorm.Session, user_id string, symbol string, freeze_amount, business_id string, behavior OpBehavior) (success bool, err error) {
	return freezeAssets(db, user_id, symbol, freeze_amount, business_id, behavior)
}

func FreezeTotalAssets(db *xorm.Session, user_id string, symbol string, business_id string, behavior OpBehavior) (success bool, err error) {
	return freezeAssets(db, user_id, symbol, "0", business_id, behavior)
}

func freezeAssets(db *xorm.Session, user_id string, symbol string, freeze_amount, business_id string, behavior OpBehavior) (success bool, err error) {

	if utils.D(freeze_amount).Cmp(utils.D("0")) < 0 {
		return false, fmt.Errorf("冻结数量必须大于等于0")
	}

	item := Assets{UserId: user_id, Symbol: symbol}
	_, err = db.Table(new(Assets)).ForUpdate().Get(&item)
	if err != nil {
		return false, err
	}

	//冻结金额为0，冻结全部可用
	if utils.D(freeze_amount).Equal(utils.D("0")) {
		freeze_amount = utils.D(item.Available).String()
	}

	item.Available = utils.D(item.Available).Sub(utils.D(freeze_amount)).String()
	item.Freeze = utils.D(item.Freeze).Add(utils.D(freeze_amount)).String()

	if utils.D(item.Available).Cmp(utils.D("0")) < 0 {
		return false, fmt.Errorf("冻结数量超出可用个数")
	}

	_, err = db.Table(new(Assets)).Where("user_id=? and symbol=?", user_id, symbol).AllCols().Update(&item)
	if err != nil {
		return false, err
	}

	//freeze log
	lg := AssetsFreeze{
		UserId:       user_id,
		Symbol:       symbol,
		Amount:       freeze_amount,
		FreezeAmount: freeze_amount,
		BusinessId:   business_id,
		Status:       FreezeStatusNew,
		Info:         string(behavior),
	}

	_, err = db.Table(&lg).Insert(&lg)
	if err != nil {
		return false, err
	}

	return true, nil
}
