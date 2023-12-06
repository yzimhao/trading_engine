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
	Ip          string     `xorm:"varchar(60)" json:"ip"`
	CreatedTime utils.Time `xorm:"timestamp created" json:"create_time"`
}

func NewAdminActivityLog(user_id int64, method string, req_uri string, body string, ip string) error {
	l := AdminActivityLog{
		UserID:     user_id,
		ActionType: method,
		ReqUri:     req_uri,
		Details:    body,
		Ip:         ip,
	}

	db := app.Database().NewSession()
	defer db.Close()

	if err := dbtables.AutoCreateTable(db, &l); err != nil {
		return err
	}
	_, err := db.Insert(&l)
	return err
}
