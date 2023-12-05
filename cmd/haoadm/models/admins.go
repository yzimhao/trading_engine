package models

import (
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app/config"
)

type AdminStatus int

const (
	AdminStatusOk AdminStatus = iota
	AdminStatusDisable
	AdminStatusLocked
)

type Admin struct {
	UserId       int64          `xorm:"pk autoincr bigint" json:"user_id"`
	Username     string         `xorm:"varchar(100) unique" json:"username"`
	Password     string         `xorm:"varchar(100)" json:"password"`
	Email        string         `xorm:"varchar(100) unique" json:"email"`
	Mobile       string         `xorm:"varchar(100) unique" json:"mobile"`
	AttemptCount int            `xorm:"default(0)" json:"attempt_count"`
	LockedAt     utils.Time     `xorm:"timestamp" json:"locked_at"`
	LoginIp      string         `xorm:"varchar(100)" json:"login_ip"`
	Role         SystemUserRole `xorm:"tinyint(1) default(0)" json:"role"`
	Status       AdminStatus    `xorm:"tinyint(1) default(0)" json:"status"`
	CreateTime   utils.Time     `xorm:"timestamp created" json:"create_time"`
	UpdateTime   utils.Time     `xorm:"timestamp updated" json:"update_time"`
}

func (a *Admin) ComparePassword(inputpasswd string) error {
	return utils.ComparePassword(config.App.Main.SecretKey, inputpasswd, a.Password)
}
