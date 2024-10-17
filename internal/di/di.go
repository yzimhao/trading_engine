package di

import (
	"context"

	"github.com/yzimhao/trading_engine/v2/app/api/handlers"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm"
	"github.com/yzimhao/trading_engine/v2/internal/services"
	"github.com/yzimhao/trading_engine/v2/internal/subscribers"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server interface {
	Start() error
	Stop() error
	Scheme() string
	Addr() string
}

func RunServer(lc fx.Lifecycle, server Server, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting server", zap.String("scheme", server.Scheme()), zap.String("addr", server.Addr()))
			go server.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping server")
			return server.Stop()
		},
	})
}

func App() *fx.App {
	return fx.New(
		fx.Provide(
			zap.NewDevelopment,
			NewViper,
			NewRedis,
			NewGinEngine,
			NewHttpServer,
			NewCache,
			NewGorm,
			NewBroker,
		),

		gorm.Module,
		handlers.Module,
		services.Module,
		subscribers.Module,
		fx.Invoke(RunServer),
	)
}
