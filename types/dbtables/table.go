package dbtables

import (
	"sync"

	"xorm.io/xorm"
)

var (
	MTables sync.Map
)

func Exist(db *xorm.Session, bean any) bool {
	table_name := db.Engine().TableName(bean)
	if _, ok := MTables.Load(table_name); ok {
		return true
	}

	exist, _ := db.IsTableExist(table_name)
	if exist {
		MTables.Store(table_name, true)
	}
	return exist

}

func CleanTable(db *xorm.Session, bean any) {
	table_name := db.Engine().TableName(bean)
	db.DropIndexes(bean)
	db.Engine().DropTables(bean)
	MTables.Delete(table_name)
}
