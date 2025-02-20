package di

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	app "github.com/yzimhao/trading_engine/v2/app"
	"github.com/yzimhao/trading_engine/v2/app/example"
	"github.com/yzimhao/trading_engine/v2/internal/modules"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database"
	"github.com/yzimhao/trading_engine/v2/internal/services"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func run(lc fx.Lifecycle, logger *zap.Logger, engine *gin.Engine, v *viper.Viper) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			v.SetDefault("port", 8080)
			engine.Run(v.GetString("port"))
			return nil
		},
		OnStop: func(ctx context.Context) error {

			return nil
		},
	})
}

func App() *fx.App {

	ctx := context.Background()
	return fx.New(
		fx.Provide(
			func() context.Context { return ctx },
			zap.NewDevelopment,
			NewViper,
			NewRedis,
			NewGin,
			NewCache,
			NewGorm,
			NewBroker,
		),

		database.Module,
		app.Module,
		example.Module,
		services.Module,
		modules.Load,
		fx.Invoke(run),
	)
}
