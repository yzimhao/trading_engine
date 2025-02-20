package app

import (
	"github.com/yzimhao/trading_engine/v2/app/api/handlers/controllers"
	"github.com/yzimhao/trading_engine/v2/app/middlewares"
	"github.com/yzimhao/trading_engine/v2/app/webws"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"app",
	fx.Provide(
		middlewares.NewAuthMiddleware,
		controllers.NewBaseController,
		controllers.NewUserAssetsController,
		controllers.NewOrderController,
		controllers.NewMarketController,
		controllers.NewUserController,
		webws.NewWsManager,
	),
	fx.Invoke(NewRoutes),
)
