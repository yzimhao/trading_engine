package di

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	app "github.com/yzimhao/trading_engine/v2/app"
	"github.com/yzimhao/trading_engine/v2/app/example"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/modules"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database"
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
	ctx, cancel := context.WithCancel(context.Background())

	app := fx.New(
		fx.Provide(
			func() context.Context {
				return ctx
			},
			zap.NewDevelopment,
			provider.NewViper,
			provider.NewRedis,
			provider.NewGin,
			provider.NewCache,
			provider.NewGorm,
			provider.NewBroker,
		),

		database.Module,
		app.Module,
		example.Module,
		modules.Load,
		fx.Invoke(run),
	)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		sig := <-c
		fmt.Printf("signal received (%v), shutting down...", sig)
		cancel()
	}()

	go func() {
		<-ctx.Done()
		fmt.Println("context cancelled, stopping fx app....")
		app.Stop(context.Background())
	}()

	return app
}
