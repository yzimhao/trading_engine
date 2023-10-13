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
			db := app.DatabaseInit()
			rc := app.RedisInit()

			base.Init(db, rc)
			assets.Init(db, rc)
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

			www.Run()
			return nil
		},
	}

	if err := appm.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
