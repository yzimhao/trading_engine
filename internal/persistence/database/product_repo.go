package database

import (
	"strings"

	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	"gorm.io/gorm"
)

type productRepo struct {
	db        *gorm.DB
	assetRepo persistence.AssetRepository
}

func NewProductRepo(datasource *gorm.DB, assetRepo persistence.AssetRepository) persistence.ProductRepository {
	return &productRepo{
		db:        datasource,
		assetRepo: assetRepo,
	}
}

func (p *productRepo) DB() *gorm.DB {
	return p.db
}

func (p *productRepo) Get(symbol string) (*entities.Product, error) {
	symbol = strings.ToLower(symbol)
	product := &entities.Product{}

	if err := p.db.Model(&entities.Product{}).Where("symbol = ?", symbol).First(&product).Error; err != nil {
		return nil, err
	}

	base, err := p.assetRepo.GetById(product.BaseId)
	if err != nil {
		return nil, err
	}
	product.Base = base

	target, err := p.assetRepo.GetById(product.TargetId)
	if err != nil {
		return nil, err
	}
	product.Target = target

	return product, nil
}
