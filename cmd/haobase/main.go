package main

import (
	"os"

	"github.com/sevlyar/go-daemon"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haobase/base/symbols"
	"github.com/yzimhao/trading_engine/cmd/haobase/clearing"
	"github.com/yzimhao/trading_engine/cmd/haobase/www"
	"github.com/yzimhao/trading_engine/utils/app"
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
			app.ConfigInit(ctx.String("config"))
			app.DatabaseInit(app.Cstring("database.driver"), app.Cstring("database.dsn"), app.Cbool("database.show_sql"))
			app.RedisInit(app.Cstring("redis.host"), app.Cstring("redis.password"), app.Cint("redis.db"))

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
					symbols.DemoData()
					assets.DemoData()
					return nil
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			app.LogsInit("haobase.run", ctx.Bool("deamon"))
			if viper.GetString("main.mode") != "prod" {
				logrus.Infof("当前运行在%s模式下，生产环境时main.mode请务必成prod", viper.GetString("main.mode"))
			}

			if ctx.Bool("deamon") {
				logrus.Info("开始守护进程")
				context, d, err := app.Deamon("haobase.pid", "")
				if err != nil {
					logrus.Fatal("创建守护进程失败, err:", err.Error())
				}
				if d != nil {
					return nil
				}

				defer func(context *daemon.Context) {
					err := context.Release()
					if err != nil {
						logrus.Printf("释放失败:%s", err.Error())
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
		logrus.Fatal(err)
	}
}

func initDemoBaseData() {
	if app.RunMode == app.ModeDemo {
		symbols.DemoData()
	}
}
