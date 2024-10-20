package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
	"github.com/yzimhao/trading_engine/v2/internal/di"
	"github.com/yzimhao/trading_engine/v2/migrations"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()

	migrateCmd := &cli.Command{
		Name:  "migrate",
		Usage: "migrate database",
		Commands: []*cli.Command{
			{
				Name:        "up",
				Description: "migrate db up",
				Action: func(ctx context.Context, c *cli.Command) error {
					return fx.New(
						fx.Provide(
							di.NewViper,
							zap.NewDevelopment,
							di.NewGorm,
						),
						fx.Invoke(migrations.MigrateUp),
					).Start(ctx)
				},
			},
			{
				Name:        "down",
				Description: "migrate db down",
				Action: func(ctx context.Context, c *cli.Command) error {
					return fx.New(
						fx.Provide(
							di.NewViper,
							zap.NewDevelopment,
							di.NewGorm,
						),
						fx.Invoke(migrations.MigrateDown),
					).Start(ctx)
				},
			},
		},
	}

	cmd := &cli.Command{
		Name: "jasmDex",
		Action: func(_ context.Context, cmd *cli.Command) error {
			app := di.App()
			app.Run()
			return nil
		},
		Commands: []*cli.Command{
			migrateCmd,
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
