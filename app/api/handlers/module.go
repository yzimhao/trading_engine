package handlers

import (
	"github.com/yzimhao/trading_engine/v2/app/api/handlers/controllers"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		controllers.NewBaseController,
		controllers.NewUserAssetsController,
		controllers.NewOrderController,
		controllers.NewMarketController,
	),
	fx.Invoke(NewRoutes),
)
