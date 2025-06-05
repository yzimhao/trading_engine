package di

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/modules"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

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

		modules.Load,
		fx.Invoke(func(lc fx.Lifecycle, logger *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error {
					logger.Info("application is starting...")
					return nil
				},
				OnStop: func(context.Context) error {
					logger.Info("application is stopping...")
					cancel()
					return nil
				},
			})
		}),
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
