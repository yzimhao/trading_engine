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

func (a *assetRepo) DB() *gorm.DB {
	return a.db
}

func (a *assetRepo) Get(symbol string) (*entities.Asset, error) {
	// todo
	var asset entities.Asset
	if err := a.db.Where("symbol = ?", symbol).First(&asset).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (a *assetRepo) GetById(id int32) (*entities.Asset, error) {
	var asset entities.Asset
	if err := a.db.Where("id = ?", id).First(&asset).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}
