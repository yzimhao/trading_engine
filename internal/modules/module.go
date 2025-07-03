package modules

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/yzimhao/trading_engine/v2/internal/modules/base"
	"github.com/yzimhao/trading_engine/v2/internal/modules/example"
	"github.com/yzimhao/trading_engine/v2/internal/modules/middlewares"
	"github.com/yzimhao/trading_engine/v2/internal/modules/notification"
	"github.com/yzimhao/trading_engine/v2/internal/modules/quote"
	"github.com/yzimhao/trading_engine/v2/internal/modules/tradingcore"
	"github.com/yzimhao/trading_engine/v2/internal/modules/usercenter"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Invoke = fx.Module(
	"modules",
	fx.Provide(
		middlewares.NewAuthMiddleware,
	),
	base.Module,
	usercenter.Module,
	tradingcore.Module,
	quote.Module,
	notification.Module,
	example.Module,
	fx.Invoke(run),
)

func run(lc fx.Lifecycle, logger *zap.Logger, engine *gin.Engine, v *viper.Viper) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			starupGinServer(v, engine)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

func starupGinServer(v *viper.Viper, engine *gin.Engine) {
	v.SetDefault("listen", "127.0.0.1")
	v.SetDefault("port", 8080)
	go engine.Run(fmt.Sprintf("%s:%d", v.GetString("listen"), v.GetInt("port")))
}
