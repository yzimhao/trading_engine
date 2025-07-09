package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/qvcloud/gopkg/version"
	"github.com/urfave/cli/v3"
	"github.com/yzimhao/trading_engine/v2/internal/di"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/migrations"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	appName      = "tradingEngine"
	appDesc      = "go开发的证券数字货币交易系统"
	appAuthor    = "yzimhao"
	appCopyright = "https://github.com/yzimhao/trading_engine"
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
							provider.NewViper,
							zap.NewDevelopment,
							provider.NewGorm,
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
							provider.NewViper,
							zap.NewDevelopment,
							provider.NewGorm,
						),
						fx.Invoke(migrations.MigrateDown),
					).Start(ctx)
				},
			},
			{
				Name:        "clean",
				Description: "clean db",
				Action: func(ctx context.Context, c *cli.Command) error {
					return fx.New(
						fx.Provide(
							provider.NewViper,
							zap.NewDevelopment,
							provider.NewGorm,
						),
						fx.Invoke(migrations.MigrateClean),
					).Start(ctx)
				},
			},
		},
	}

	versionCmd := &cli.Command{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Display version info.",
		Action: func(ctx context.Context, c *cli.Command) error {
			version.ShowVersion()
			return nil
		},
	}

	cmd := &cli.Command{
		Name:        appName,
		Description: appDesc,
		Usage:       "",
		Authors:     []any{appAuthor},
		Copyright:   appCopyright,
		Action: func(_ context.Context, cmd *cli.Command) error {
			app := di.App()
			app.Run()
			return nil
		},
		Commands: []*cli.Command{
			migrateCmd,
			versionCmd,
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
