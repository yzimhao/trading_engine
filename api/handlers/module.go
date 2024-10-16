package handlers

import (
	"github.com/yzimhao/trading_engine/v2/api/handlers/controllers"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		controllers.NewUserAssetsController,
	),
	fx.Invoke(NewRoutesContext),
)
