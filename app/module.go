package app

import (
	"github.com/yzimhao/trading_engine/v2/app/api/handlers/controllers"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"app",
	fx.Provide(
		controllers.NewBaseController,
		controllers.NewUserAssetsController,
		controllers.NewOrderController,
		controllers.NewMarketController,
		controllers.NewUserController,
	),
	fx.Invoke(NewRoutes),
)
