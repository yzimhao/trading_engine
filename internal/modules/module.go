package modules

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/yzimhao/trading_engine/v2/app"
	"github.com/yzimhao/trading_engine/v2/app/example"
	"github.com/yzimhao/trading_engine/v2/internal/modules/matching"
	"github.com/yzimhao/trading_engine/v2/internal/modules/quote"
	"github.com/yzimhao/trading_engine/v2/internal/modules/settlement"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Load = fx.Module(
	"modules",
	app.Module,
	example.Module,

	settlement.Module,
	matching.Module,
	quote.Module,
	fx.Invoke(run),
)

func run(lc fx.Lifecycle, logger *zap.Logger, engine *gin.Engine, v *viper.Viper) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			v.SetDefault("listen", "127.0.0.1")
			v.SetDefault("port", 8080)
			go engine.Run(fmt.Sprintf("%s:%d", v.GetString("listen"), v.GetInt("port")))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
