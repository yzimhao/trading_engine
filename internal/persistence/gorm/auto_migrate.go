package gorm

import (
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
	"gorm.io/gorm"
)

func autoMigrate(db *gorm.DB) error {

	return nil

	return db.AutoMigrate(
		&entities.Assets{},
		&entities.AssetsLog{},
		&entities.AssetsFreeze{},
	)

}
