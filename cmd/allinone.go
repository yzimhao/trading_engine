package main

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sevlyar/go-daemon"
	"github.com/urfave/cli/v2"
	admview "github.com/yzimhao/trading_engine/cmd/haoadm/view"
	base "github.com/yzimhao/trading_engine/cmd/haobase/run"
	"github.com/yzimhao/trading_engine/cmd/haomatch/matching"
	"github.com/yzimhao/trading_engine/cmd/haoquote/quote"
	"github.com/yzimhao/trading_engine/cmd/haosettle/settle"
	"github.com/yzimhao/trading_engine/config"
	"github.com/yzimhao/trading_engine/utils/app"
	"github.com/yzimhao/trading_engine/utils/app/keepalive"

	_ "net/http/pprof"
)

func main() {

	appm := &cli.App{
		Name:      "haotrader",
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
			if ctx.Bool("deamon") {
				context, d, err := app.Deamon(ctx.String("pid"), "")
				if err != nil {
					app.Logger.Fatal("创建守护进程失败, err:", err.Error())
				}
				if d != nil {
					app.Logger.Printf("这是在父进程的标志")
					return nil
				}

				defer func(context *daemon.Context) {
					err := context.Release()
					if err != nil {
						app.Logger.Printf("释放失败:%s", err.Error())
					}
					app.Logger.Printf("释放成功!!!")
				}(context)
			}

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
	debug_pprof()
	go haobase()
	//todo 加入context，优化掉这个sleep
	time.Sleep(time.Second)
	go haomatch()
	go haosettle()
	go haoquote()
	haoadm()
}

func debug_pprof() {
	if config.App.Main.Mode != app.ModeProd.String() {
		go func() {
			app.Logger.Info(http.ListenAndServe("0.0.0.0:26060", nil))
		}()
	}
}

func haobase() {
	base.Run()
}

func haomatch() {
	matching.Run()
}

func haosettle() {
	settle.Run()
}

func haoquote() {
	quote.Run()
}

func haoadm() {
	admview.Run()
}
