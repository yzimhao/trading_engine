package dbtables

import (
	"sync"

	"xorm.io/xorm"
)

var (
	MTables sync.Map
)

func Exist(db *xorm.Session, table_name string) bool {
	if _, ok := MTables.Load(table_name); ok {
		return true
	}

	exist, _ := db.IsTableExist(table_name)
	if exist {
		MTables.Store(table_name, true)
	}
	return exist

}
