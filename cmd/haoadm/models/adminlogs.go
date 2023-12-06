package models

import (
	"github.com/yzimhao/trading_engine/types/dbtables"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

type AdminActivityLog struct {
	Id          int64      `xorm:"pk autoincr bigint" json:"id"`
	UserID      int64      `xorm:"user_id" json:"user_id"`
	ActionType  string     `xorm:"varchar(30)" json:"action_type"`
	ReqUri      string     `xorm:"varchar(300)" json:"req_uri"`
	Details     string     `xorm:"text" json:"details"`
	CreatedTime utils.Time `xorm:"timestamp created" json:"create_time"`
}

func NewAdminActivityLog(user_id int64, method string, req_uri string, body string) {
	l := AdminActivityLog{
		UserID:     user_id,
		ActionType: method,
		ReqUri:     req_uri,
		Details:    body,
	}

	db := app.Database().NewSession()
	defer db.Close()

	dbtables.AutoCreateTable(db, &l)
	db.Insert(&l)
}
