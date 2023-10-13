package assets

import (
	"fmt"

	"xorm.io/xorm"
)

func UnfreezeAssets(db *xorm.Session, user_id string, business_id, unfreeze_amount string) (success bool, err error) {
	return unfreezeAssets(db, user_id, business_id, unfreeze_amount)
}

func UnfreezeAllAssets(db *xorm.Session, user_id string, business_id string) (success bool, err error) {
	return unfreezeAssets(db, user_id, business_id, "0")
}

func unfreezeAssets(db *xorm.Session, user_id string, business_id, unfreeze_amount string) (success bool, err error) {
	if check_number_lt_zero(unfreeze_amount) {
		return false, fmt.Errorf("unfreeze amount should be >= 0")
	}

	row := assetsFreeze{UserId: user_id, BusinessId: business_id}

	has, err := db.Table(new(assetsFreeze)).Where("business_id=?", business_id).Get(&row)
	if err != nil {
		return false, err
	}

	if !has {
		return false, fmt.Errorf("not found business_id")
	}

	if row.Status == FreezeStatusDone {
		return false, fmt.Errorf("repeat unfreeze")
	}

	if d(unfreeze_amount).Equal(d("0")) {
		unfreeze_amount = row.FreezeAmount
	}

	row.FreezeAmount = number_sub(row.FreezeAmount, unfreeze_amount)

	if check_number_lt_zero(row.FreezeAmount) {
		return false, fmt.Errorf("unfreeze amount must lt freeze amount")
	}

	if check_number_eq_zero(row.FreezeAmount) {
		row.Status = FreezeStatusDone
	}

	_, err = db.Table(new(assetsFreeze)).Where("business_id=?", business_id).AllCols().Update(&row)
	if err != nil {
		return false, err
	}

	//解冻资产为可用
	assets := Assets{UserId: user_id, Symbol: row.Symbol}
	_, err = db.Table(new(Assets)).Where("user_id=? and symbol=?", user_id, row.Symbol).ForUpdate().Get(&assets)
	if err != nil {
		return false, err
	}
	assets.Available = number_add(assets.Available, unfreeze_amount)
	assets.Freeze = number_sub(assets.Freeze, unfreeze_amount)

	if check_number_lt_zero(assets.Freeze) {
		return false, fmt.Errorf("freeze amount some wrong")
	}

	_, err = db.Table(new(Assets)).Where("user_id=? and symbol=?", user_id, row.Symbol).Update(&assets)
	if err != nil {
		return false, err
	}

	return true, nil
}
