package database

import (
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	"gorm.io/gorm"
)

type assetRepo struct {
	db *gorm.DB
}

func NewAssetRepo(datasource *gorm.DB) persistence.AssetRepository {

	return &assetRepo{
		db: datasource,
	}
}

func (v *assetRepo) Get(symbol string) (*entities.Asset, error) {
	return nil, nil
}
