package main

import (
	"os"

	"github.com/sevlyar/go-daemon"
	"github.com/urfave/cli/v2"
	"github.com/yzimhao/trading_engine/cmd/haoquote/quote"
	"github.com/yzimhao/trading_engine/utils/app"
	"github.com/yzimhao/trading_engine/utils/app/config"
	"github.com/yzimhao/trading_engine/utils/app/keepalive"
)

func main() {
	appm := &cli.App{
		Name:      "haoquote",
		UsageText: "Issues: https://github.com/yzimhao/trading_engine/issues",
		Usage:     "交易行情系统",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Value: "./config.toml", Aliases: []string{"c"}},
			&cli.StringFlag{Name: "pid", Value: "./run/haoquote.pid"},
			&cli.BoolFlag{Name: "deamon", Value: false, Aliases: []string{"d"}},
		},

		Before: func(ctx *cli.Context) error {
			app.ConfigInit(ctx.String("config"), ctx.Bool("deamon"))
			app.DatabaseInit(config.App.Database.Driver, config.App.Database.DSN, config.App.Database.ShowSQL, config.App.Database.Prefix)
			app.RedisInit(config.App.Redis.Host, config.App.Redis.Password, config.App.Redis.DB)
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
		},
		Action: func(ctx *cli.Context) error {
			if ctx.Bool("deamon") {

				context, d, err := app.Deamon(ctx.String("pid"), "")
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
			keepalive.NewKeepalive(app.RedisPool(), ctx.App.Name, app.Version, 5)
			quote.Run()
			return nil
		},
	}

	if err := appm.Run(os.Args); err != nil {
		app.Logger.Fatal(err)
	}
}
