package database

import (
	"github.com/duolacloud/crud-core/datasource"
	order_repo "github.com/yzimhao/trading_engine/v2/internal/persistence/database/order"
	"go.uber.org/fx"
	_gorm "gorm.io/gorm"
)

var Module = fx.Module(
	"database",
	fx.Provide(
		datasource.NewDataSource[_gorm.DB],
		NewAssetRepo,
		NewUserAssetRepo,
		NewProductRepo,
		order_repo.NewOrderRepo,
		NewKlineRepo,
		NewTradeRecordRepo,
	),

	fx.Invoke(autoMigrate),
)
