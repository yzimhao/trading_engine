package assets

import (
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

type AssetsHealthStatus struct {
	GlobalTotal string `json:"global_total"`
	UserTotal   string `json:"user_total"`
	Health      bool   `json:"health"`
}

func AssetsCheck() AssetsHealthStatus {
	db := app.Database().NewSession()
	defer db.Close()

	var global_total string
	db.Table(new(Assets)).Select("sum(total)").Get(&global_total)

	var user_total string
	db.Table(new(Assets)).Where("user_id !=?", "root").Select("sum(total)").Get(&user_total)

	return AssetsHealthStatus{
		GlobalTotal: global_total,
		UserTotal:   user_total,
		Health: func() bool {
			if utils.D(global_total).Equal(utils.D("0")) {
				return true
			}
			return false
		}(),
	}
}
