package persistence

import (
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
)

type ProductRepository interface {
	Find(symbol string) ([]*entities.Product, error)
}
