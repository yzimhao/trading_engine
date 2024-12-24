package app

import (
	"github.com/yzimhao/trading_engine/v2/app/api/handlers/controllers"
	"github.com/yzimhao/trading_engine/v2/app/middlewares"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		middlewares.NewAuthMiddleware,
		controllers.NewBaseController,
		controllers.NewUserAssetsController,
		controllers.NewOrderController,
		controllers.NewMarketController,
		controllers.NewUserController,
	),
	fx.Invoke(NewRoutes),
)
