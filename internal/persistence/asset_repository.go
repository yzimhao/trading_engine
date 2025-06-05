package persistence

import "github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"

type AssetRepository interface {
	Get(symbol string) (*entities.Asset, error)
}
