package persistence

import (
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	"gorm.io/gorm"
)

type ProductRepository interface {
	DB() *gorm.DB
	Get(symbol string) (*entities.Product, error)
}
