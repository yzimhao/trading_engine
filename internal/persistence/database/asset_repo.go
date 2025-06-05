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
	return nil, nil
}

func (a *assetRepo) GetById(id int32) (*entities.Asset, error) {
	// todo
	return nil, nil
}
