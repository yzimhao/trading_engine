package gorm

import (
	"context"

	"github.com/duolacloud/crud-core/types"
	models_types "github.com/yzimhao/trading_engine/v2/internal/models/types"
	models_variety "github.com/yzimhao/trading_engine/v2/internal/models/variety"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type inContext struct {
	fx.In
	Db               *gorm.DB
	Logger           *zap.Logger
	VarietyRepo      persistence.VarietyRepository
	TradeVarietyRepo persistence.TradeVarietyRepository
}

func autoMigrate(in inContext) error {

	// auto migrate
	err := in.Db.AutoMigrate(
		&entities.Asset{},
		&entities.AssetLog{},
		&entities.AssetFreeze{},
		&entities.Variety{},
		&entities.TradeVariety{},
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
	varieties, err := initVariety(ctx, in.VarietyRepo)
	if err != nil {
		return err
	}

	err = initTradeVariety(ctx, in.TradeVarietyRepo, varieties)
	if err != nil {
		return err
	}

	return nil
}

func initVariety(ctx context.Context, varietyRepo persistence.VarietyRepository) ([]*models_variety.Variety, error) {

	usdt, _ := varietyRepo.QueryOne(ctx, map[string]any{
		"symbol": map[string]any{
			"eq": "usdt",
		},
	})

	btc, _ := varietyRepo.QueryOne(ctx, map[string]any{
		"symbol": map[string]any{
			"eq": "btc",
		},
	})

	if usdt != nil && btc != nil {
		return []*models_variety.Variety{usdt, btc}, nil
	}

	opts := []types.CreateManyOption{
		types.WithCreateBatchSize(2),
	}
	varieties, err := varietyRepo.CreateMany(ctx, []*models_variety.CreateVariety{
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
	}, opts...)
	if err != nil {
		return nil, err
	}
	return varieties, nil
}

func initTradeVariety(ctx context.Context, tradeVarietyRepo persistence.TradeVarietyRepository, varieties []*models_variety.Variety) error {

	btcusdt, _ := tradeVarietyRepo.QueryOne(ctx, map[string]any{
		"symbol": map[string]any{
			"eq": "btcusdt",
		},
	})
	if btcusdt != nil {
		return nil
	}
	_, err := tradeVarietyRepo.Create(ctx, &models_variety.CreateTradeVariety{
		Symbol:         "btcusdt", //统一用小写
		Name:           "BTCUSDT",
		BaseId:         varieties[0].ID,
		TargetId:       varieties[1].ID,
		PriceDecimals:  2,
		QtyDecimals:    6,
		AllowMinQty:    "0.0001",
		AllowMinAmount: "1.00",
		AllowMaxAmount: "0",
		FeeRate:        "0.005",
		Status:         models_types.StatusEnabled,
	}, nil)
	if err != nil {
		return err
	}
	return nil
}
