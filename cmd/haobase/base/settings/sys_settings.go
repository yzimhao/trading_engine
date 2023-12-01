package settings

import (
	"github.com/yzimhao/trading_engine/types"
	"github.com/yzimhao/trading_engine/utils"
)

type ParamType string

const (
	ParamTypeInt     ParamType = "int"
	ParamTypeString  ParamType = "string"
	ParamTypeDecimal ParamType = "decimal"
	ParamTypeJson    ParamType = "json"
	ParamTypeDate    ParamType = "date"
	ParamTypeTime    ParamType = "time"
)

type Setting struct {
	Id          int          `xorm:"pk autoincr int" json:"id"`
	Code        string       `xorm:"varchar(100) notnull unique(code)" json:"code"`
	Name        string       `xorm:"varchar(250)" json:"name"`
	Type        ParamType    `xorm:"varchar(20)" json:"type"`
	Value       string       `xorm:"text" json:"value"`
	Status      types.Status `xorm:"default(0) notnull" json:"status"`
	EffectiveSt utils.Time   `xorm:"timestamp" json:"effective_st"`
	EffectiveEt utils.Time   `xorm:"timestamp" json:"effective_et"`
	CreateTime  utils.Time   `xorm:"timestamp created" json:"create_time"`
	UpdateTime  utils.Time   `xorm:"timestamp updated" json:"update_time"`
}
