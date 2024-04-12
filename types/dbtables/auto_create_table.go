package dbtables

import (
	"xorm.io/xorm"
)

func AutoCreateTable(db *xorm.Session, bean any) error {
	table := db.Engine().TableName(bean)
	if !Exist(db, table) {
		if err := db.Engine().Sync2(bean); err != nil {
			return err
		}
	}
	return nil
}
