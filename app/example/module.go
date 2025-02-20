package example

import (
	"github.com/yzimhao/trading_engine/v2/app/example/controllers"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"example",
	fx.Invoke(controllers.NewExampleController),
)
