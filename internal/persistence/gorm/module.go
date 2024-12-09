package gorm

import (
	"github.com/duolacloud/crud-core/datasource"
	order_repo "github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/order"
	"go.uber.org/fx"
	_gorm "gorm.io/gorm"
)

var Module = fx.Options(
	fx.Provide(
		datasource.NewDataSource[_gorm.DB],
		NewAssetRepo,
		NewVarietyRepo,
		NewTradeVarietyRepo,
		order_repo.NewOrderRepo,
		NewKlineRepo,
		NewTradeLogRepo,
	),

	fx.Invoke(autoMigrate),
)
