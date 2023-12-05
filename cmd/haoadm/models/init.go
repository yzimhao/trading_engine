package models

import (
	"github.com/yzimhao/trading_engine/types/dbtables"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
	"github.com/yzimhao/trading_engine/utils/app/config"
	"xorm.io/xorm"
)

func Init() {
	db := app.Database().NewSession()
	defer db.Close()

	if !dbtables.Exist(db, &Admin{}) {
		init_default_admin(db)
	}
}

func init_default_admin(db *xorm.Session) {
	dbtables.AutoCreateTable(db, &Admin{})

	a := Admin{
		Username: "admin",
		Password: Passwd("admin"),
		Role:     SystemUserRoleRoot,
		Status:   AdminStatusOk,
	}
	app.Logger.Infof("default admin: %v", a)
	db.Insert(&a)
}

func Passwd(passwd string) string {
	hash, err := utils.Password(config.App.Main.SecretKey, passwd)
	if err != nil {
		app.Logger.Errorf("utils.Password: %s", err.Error())
	}
	return hash
}
