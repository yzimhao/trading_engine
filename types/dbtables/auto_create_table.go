package dbtables

import (
	"github.com/yzimhao/trading_engine/utils/app"
	"xorm.io/xorm"
)

func AutoCreateTable(db *xorm.Session, bean any) error {
	table := db.Engine().TableName(bean)

	app.Logger.Debugf("table_name: %s %#v", table, bean)

	if !Exist(db, table) {
		err := db.CreateTable(bean)
		if err != nil {
			return err
		}

		err = db.CreateIndexes(bean)
		if err != nil {
			return err
		}

		err = db.CreateUniques(bean)
		if err != nil {
			return err
		}
	}
	return nil
}
