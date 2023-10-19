package main

import (
	"os"

	"github.com/sevlyar/go-daemon"

	"github.com/urfave/cli/v2"
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haobase/base/varieties"
	"github.com/yzimhao/trading_engine/cmd/haobase/clearing"
	"github.com/yzimhao/trading_engine/cmd/haobase/www"
	"github.com/yzimhao/trading_engine/utils/app"
	"github.com/yzimhao/trading_engine/utils/app/config"
)

func main() {
	appm := &cli.App{
		Name:      "haobase",
		UsageText: "Issues: https://github.com/yzimhao/trading_engine/issues",
		Usage:     "交易所基础模块",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Value: "./config.toml", Aliases: []string{"c"}},
			&cli.BoolFlag{Name: "deamon", Value: false, Aliases: []string{"d"}},
		},

		Before: func(ctx *cli.Context) error {
			app.ConfigInit(ctx.String("config"), ctx.Bool("deamon"))
			app.DatabaseInit(config.App.Database.Driver, config.App.Database.DSN, config.App.Database.ShowSQL, config.App.Database.Prefix)
			app.RedisInit(config.App.Redis.Host, config.App.Redis.Password, config.App.Redis.DB)

			base.Init()
			assets.Init()
			return nil
		},

		Commands: []*cli.Command{
			{
				Name: "version",
				Action: func(ctx *cli.Context) error {
					app.ShowVersion()
					return nil
				},
			},
			{
				Name: "demo-data",
				Action: func(ctx *cli.Context) error {
					varieties.DemoData()
					assets.DemoData()
					return nil
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.Bool("deamon") {

				context, d, err := app.Deamon("haobase.pid", "")
				if err != nil {
					app.Logger.Fatal("创建守护进程失败, err:", err.Error())
				}
				if d != nil {
					return nil
				}

				defer func(context *daemon.Context) {
					err := context.Release()
					if err != nil {
						app.Logger.Printf("释放失败:%s", err.Error())
					}
				}(context)

			}

			initDemoBaseData()
			clearing.Run()
			www.Run()
			return nil
		},
	}

	if err := appm.Run(os.Args); err != nil {
		app.Logger.Fatal(err)
	}
}

func initDemoBaseData() {
	if config.App.Main.Mode == config.ModeDemo {
		varieties.DemoData()
	}
}
