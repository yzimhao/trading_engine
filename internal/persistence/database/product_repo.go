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

func (v *productRepo) Find(symbol string) ([]*entities.Product, error) {
	symbol = strings.ToLower(symbol)

	if err := v.db.Model(&models_variety.TradeVariety{}).Where("symbol = ?", symbol).First(&tradeVariety).Error; err != nil {
		return nil, err
	}

	tradeVariety.BaseVariety, err = v.varietyRepo.QueryOne(ctx, map[string]any{
		"id": map[string]any{
			"eq": tradeVariety.BaseId,
		},
	})
	if err != nil {
		return nil, err
	}

	tradeVariety.TargetVariety, err = v.varietyRepo.QueryOne(ctx, map[string]any{
		"id": map[string]any{
			"eq": tradeVariety.TargetId,
		},
	})
	if err != nil {
		return nil, err
	}

	return tradeVariety, nil
}
