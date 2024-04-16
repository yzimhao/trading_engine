package main

import (
	"os"
	"path/filepath"

	"github.com/sevlyar/go-daemon"

	"github.com/urfave/cli/v2"
	"github.com/yzimhao/trading_engine/cmd/haosettle/settle"
	"github.com/yzimhao/trading_engine/config"
	"github.com/yzimhao/trading_engine/utils/app"
	"github.com/yzimhao/trading_engine/utils/app/keepalive"
)

func main() {
	appm := &cli.App{
		Name:      "haosettle",
		UsageText: "Issues: https://github.com/yzimhao/trading_engine/issues",
		Usage:     "",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Value: "./config.toml", Aliases: []string{"c"}},
			&cli.StringFlag{Name: "pid", Value: "./run/haosettle.pid"},
			&cli.BoolFlag{Name: "deamon", Value: false, Aliases: []string{"d"}},
		},

		Before: func(ctx *cli.Context) error {
			app.ConfigInit(ctx.String("config"), config.App)
			app.LogsInit(filepath.Base(os.Args[0]), config.App.Main.LogPath, config.App.Main.LogLevel, !ctx.Bool("deamon"))
			app.TimeZoneInit(config.App.Main.TimeZone)
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
			{
				Name: "settle",
				Action: func(ctx *cli.Context) error {
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
			settle.Run()
			return nil
		},
	}

	if err := appm.Run(os.Args); err != nil {
		app.Logger.Fatal(err)
	}
}
