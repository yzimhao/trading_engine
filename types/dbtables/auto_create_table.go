package dbtables

import (
	"github.com/yzimhao/trading_engine/utils/app"
	"xorm.io/xorm"
)

func AutoCreateTable(db *xorm.Session, bean any) error {
	table := db.Engine().TableName(bean)

	app.Logger.Debugf("auto create table_name: %s", table)

	if !Exist(db, table) {
		if err := db.Engine().Sync2(bean); err != nil {
			return err
		}

		// err := db.CreateTable(bean)
		// if err != nil {
		// 	return err
		// }

		// err = db.CreateIndexes(bean)
		// if err != nil {
		// 	return err
		// }

		// err = db.CreateUniques(bean)
		// if err != nil {
		// 	return err
		// }
	}
	return nil
}
