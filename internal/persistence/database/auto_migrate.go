package database

import (
	"context"

	models_types "github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type inContext struct {
	fx.In
	Db          *gorm.DB
	Logger      *zap.Logger
	AssetRepo   persistence.AssetRepository
	ProductRepo persistence.ProductRepository
}

func autoMigrate(in inContext) error {

	// auto migrate
	err := in.Db.AutoMigrate(
		&entities.UserAsset{},
		&entities.UserAssetLog{},
		&entities.UserAssetFreeze{},
		&entities.Asset{},
		&entities.Product{},
	)
	if err != nil {
		in.Logger.Error("auto migrate error", zap.Error(err))
		return err
	}

	//init data
	err = initData(context.Background(), in)
	if err != nil {
		in.Logger.Error("init data error", zap.Error(err))
		return err
	}

	return nil
}

func initData(ctx context.Context, in inContext) error {
	assets, err := initAsset(ctx, in.AssetRepo)
	if err != nil {
		return err
	}

	err = initProduct(ctx, in.ProductRepo, assets)
	if err != nil {
		return err
	}

	return nil
}

func initAsset(ctx context.Context, assetRepo persistence.AssetRepository) ([]*entities.Asset, error) {

	usdt, _ := assetRepo.Get("usdt")
	btc, _ := assetRepo.Get("btc")

	if usdt != nil && btc != nil {
		return []*entities.Asset{usdt, btc}, nil
	}

	assets := []*entities.Asset{
		{
			Symbol:       "usdt",
			Name:         "USDT",
			ShowDecimals: 4,
			MinDecimals:  6,
			IsBase:       true,
			Status:       models_types.StatusEnabled,
		},
		{
			Symbol:       "btc",
			Name:         "Bitcoin",
			ShowDecimals: 4,
			MinDecimals:  8,
			Status:       models_types.StatusEnabled,
		},
	}

	if err := assetRepo.DB().CreateInBatches(assets, 2).Error; err != nil {
		return nil, err
	}
	return assets, nil
}

func initProduct(ctx context.Context, productRepo persistence.ProductRepository, assets []*entities.Asset) error {

	btcusdt, _ := productRepo.Get("btcusdt")
	if btcusdt != nil {
		return nil
	}
	conn := productRepo.DB().Create(&entities.Product{
		Symbol:         "btcusdt", //统一用小写
		Name:           "BTCUSDT",
		BaseId:         assets[0].ID,
		TargetId:       assets[1].ID,
		PriceDecimals:  2,
		QtyDecimals:    6,
		AllowMinQty:    "0.0001",
		AllowMinAmount: "1.00",
		AllowMaxAmount: "0",
		FeeRate:        "0.005",
		Status:         models_types.StatusEnabled,
	})
	if err := conn.Error; err != nil {
		return err
	}
	return nil
}
