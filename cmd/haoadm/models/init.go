package models

import (
	"github.com/yzimhao/trading_engine/config"
	"github.com/yzimhao/trading_engine/types/dbtables"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
	"xorm.io/xorm"
)

func Init() {
	db := app.Database().NewSession()
	defer db.Close()

	if !dbtables.Exist(db, &Adminuser{}) {
		init_default_admin(db)
	}
}

func init_default_admin(db *xorm.Session) {
	dbtables.AutoCreateTable(db, &Adminuser{})

	a := Adminuser{
		Username: "admin",
		Password: Passwd("admin"),
		Email:    "admin@admin.com",
		Mobile:   "+1(502) 414-5497",
		Role:     AdminRoleSuper,
		Status:   AdminuserStatusNormal,
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
