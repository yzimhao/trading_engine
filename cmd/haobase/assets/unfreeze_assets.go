package assets

import (
	"fmt"

	"github.com/yzimhao/trading_engine/utils"
	"xorm.io/xorm"
)

func UnfreezeAssets(db *xorm.Session, user_id string, business_id, unfreeze_amount string) (success bool, err error) {
	return unfreezeAssets(db, user_id, business_id, unfreeze_amount)
}

func UnfreezeAllAssets(db *xorm.Session, user_id string, business_id string) (success bool, err error) {
	return unfreezeAssets(db, user_id, business_id, "0")
}

func unfreezeAssets(db *xorm.Session, user_id string, business_id, unfreeze_amount string) (success bool, err error) {

	if utils.D(unfreeze_amount).Cmp(utils.D("0")) >= 0 {
		return false, fmt.Errorf("解冻金额必须大于等于0")
	}

	row := assetsFreeze{UserId: user_id, BusinessId: business_id}

	has, err := db.Table(new(assetsFreeze)).Where("business_id=?", business_id).Get(&row)
	if err != nil {
		return false, err
	}

	if !has {
		return false, fmt.Errorf("未找到冻结记录")
	}

	if row.Status == FreezeStatusDone {
		return false, fmt.Errorf("冻结记录已经解冻")
	}

	//解冻金额为0，则解冻全部
	if utils.D(unfreeze_amount).Equal(utils.D("0")) {
		unfreeze_amount = row.FreezeAmount
	}

	freeze_amount := utils.D(row.FreezeAmount).Sub(utils.D(unfreeze_amount))

	if freeze_amount.Cmp(utils.D("0")) < 0 {
		return false, fmt.Errorf("数据错误，解冻后金额为负数")
	}

	if freeze_amount.Equal(utils.D("0")) {
		row.Status = FreezeStatusDone
	}

	row.FreezeAmount = freeze_amount.String()
	_, err = db.Table(new(assetsFreeze)).Where("business_id=?", business_id).AllCols().Update(&row)
	if err != nil {
		return false, err
	}

	//解冻资产为可用
	assets := Assets{UserId: user_id, Symbol: row.Symbol}
	_, err = db.Table(new(Assets)).Where("user_id=? and symbol=?", user_id, row.Symbol).Get(&assets)
	if err != nil {
		return false, err
	}
	assets.Available = utils.D(assets.Available).Add(utils.D(unfreeze_amount)).String()
	assets.Freeze = utils.D(assets.Freeze).Sub(utils.D(unfreeze_amount)).String()

	if utils.D(assets.Freeze).Cmp(utils.D("0")) < 0 {
		return false, fmt.Errorf("数据出错，冻结金额为负数")
	}

	_, err = db.Table(new(Assets)).Where("user_id=? and symbol=?", user_id, row.Symbol).Update(&assets)
	if err != nil {
		return false, err
	}
	return true, nil
}
