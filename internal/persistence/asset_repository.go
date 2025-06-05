package persistence

import (
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	"gorm.io/gorm"
)

type AssetRepository interface {
	DB() *gorm.DB
	Get(symbol string) (*entities.Asset, error)
	GetById(id int32) (*entities.Asset, error)
}
