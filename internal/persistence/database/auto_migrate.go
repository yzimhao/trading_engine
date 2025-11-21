package database

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	"github.com/yzimhao/trading_engine/v2/internal/types"
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

	// 按模型逐表处理迁移：
	// - 如果表不存在：使用 AutoMigrate 创建表（首次创建包含所有字段）
	// - 如果表已存在：仅为缺失字段执行 AddColumn，跳过对已存在相同字段的任何更改
	var err error

	models := []interface{}{
		&entities.UserAsset{},
		&entities.UserAssetLog{},
		&entities.UserAssetFreeze{},
		&entities.Asset{},
		&entities.Product{},
	}

	for _, m := range models {
		hasTable := in.Db.Migrator().HasTable(m)
		if !hasTable {
			if err := in.Db.AutoMigrate(m); err != nil {
				in.Logger.Error("auto migrate create table error", zap.Error(err))
				return err
			}
			continue
		}

		// 表已存在，仅添加缺失字段
		if err := addMissingColumns(in.Db, m); err != nil {
			in.Logger.Error("add missing columns error", zap.Error(err))
			return err
		}
	}

	//init data
	err = initData(context.Background(), in)
	if err != nil {
		in.Logger.Error("init data error", zap.Error(err))
		return err
	}

	return nil
}

// addMissingColumns 仅为模型中在数据库中不存在的字段添加列；不修改已有字段
func addMissingColumns(db *gorm.DB, model interface{}) error {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return err
	}

	for _, field := range stmt.Schema.Fields {
		// 使用字段的 DB 列名（DBName）进行存在性判断；有些嵌入字段或未映射字段可能没有 DBName，跳过这些字段
		colName := field.DBName
		if colName == "" {
			// 跳过没有映射到数据库列的字段
			continue
		}

		// GORM 的 Migrator.HasColumn 在不同版本对参数的接受形式不完全一致，先尝试用列名判断，再尝试用字段名
		has := false
		if db.Migrator().HasColumn(model, colName) {
			has = true
		} else if db.Migrator().HasColumn(model, field.Name) {
			has = true
		}

		if !has {
			// 使用字段名调用 AddColumn（GORM 期望字段名或列名，两者通常都可），但避免传空字符串
			if field.Name == "" {
				return fmt.Errorf("cannot add column for model: empty field name for column %s", colName)
			}
			if err := db.Migrator().AddColumn(model, field.Name); err != nil {
				return fmt.Errorf("failed to add column %s: %w", colName, err)
			}
		}
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
			Status:       types.StatusEnabled,
		},
		{
			Symbol:       "btc",
			Name:         "Bitcoin",
			ShowDecimals: 4,
			MinDecimals:  8,
			Status:       types.StatusEnabled,
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
		AllowMinQty:    decimal.NewFromFloat(0.0001),
		AllowMinAmount: decimal.NewFromFloat(1.0),
		AllowMaxAmount: decimal.Zero,
		FeeRate:        decimal.NewFromFloat(0.005),
		Status:         types.StatusEnabled,
	})
	if err := conn.Error; err != nil {
		return err
	}
	return nil
}
