package main

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
	admview "github.com/yzimhao/trading_engine/cmd/haoadm/view"
	"github.com/yzimhao/trading_engine/cmd/haobase/core"
	"github.com/yzimhao/trading_engine/cmd/haomatch/matching"
	"github.com/yzimhao/trading_engine/cmd/haoquote/quote"
	"github.com/yzimhao/trading_engine/config"
	"github.com/yzimhao/trading_engine/utils/app"
	"github.com/yzimhao/trading_engine/utils/app/keepalive"
)

func main() {

	appm := &cli.App{
		Name:      "allinone",
		UsageText: "Issues: https://github.com/yzimhao/trading_engine/issues",
		Usage:     "交易软件",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Value: "./config.toml", Aliases: []string{"c"}},
			&cli.StringFlag{Name: "pid", Value: "./run/haoadm.pid"},
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
		},
		Action: func(ctx *cli.Context) error {
			keepalive.NewKeepalive(app.RedisPool(), ctx.App.Name, app.Version, 5)
			start()
			return nil
		},
	}

	if err := appm.Run(os.Args); err != nil {
		app.Logger.Fatal(err)
	}
}

func start() {
	go haobase()
	go haomatch()
	go haoquote()
	haoadm()
}

func haobase() {
	core.Run()
}

func haomatch() {
	matching.Run()
}

func haoquote() {
	quote.Run()
}

func haoadm() {
	admview.Run()
}
