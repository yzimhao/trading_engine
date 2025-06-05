package persistence

import (
	"github.com/shopspring/decimal"
	models "github.com/yzimhao/trading_engine/v2/internal/models/asset"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	"gorm.io/gorm"
)

type UserAssetRepository interface {
	QueryUserAsset(userId string, symbol string) (*entities.UserAsset, error)
	QueryUserAssets(userId string, symbols ...string) ([]*entities.UserAsset, error)

	Despoit(transId, userId, symbol string, amount decimal.Decimal) error
	Withdraw(transId, userId, symbol string, amount decimal.Decimal) error
	Transfer(transId, from, to, symbol string, amount decimal.Decimal) error
	TransferWithTx(tx *gorm.DB, transId, from, to, symbol string, amount decimal.Decimal) error
	Freeze(tx *gorm.DB, transId, userId, symbol string, amount decimal.Decimal) (*entities.UserAssetFreeze, error)
	UnFreeze(tx *gorm.DB, transId, userId, symbol string, amount decimal.Decimal) error
	QueryFreeze(filter map[string]any) (assetFreezes []*models.AssetFreeze, err error)
}
