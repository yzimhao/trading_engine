package main

import (
	"os"

	"github.com/sevlyar/go-daemon"
	"github.com/urfave/cli/v2"
	"github.com/yzimhao/trading_engine/cmd/haomatch/matching"
	"github.com/yzimhao/trading_engine/utils/app"
	"github.com/yzimhao/trading_engine/utils/app/config"
)

func main() {
	appm := &cli.App{
		Name:      "haomatch",
		UsageText: "Issues: https://github.com/yzimhao/trading_engine/issues",
		Usage:     "交易撮合引擎",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Value: "./config.toml", Aliases: []string{"c"}},
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

			{
				Name: "test",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "symbol", Required: true, Aliases: []string{"s"}, DefaultText: "", Usage: "usdjpy"},
					&cli.StringFlag{Name: "side", Required: true, Value: "ask", Usage: "ask/bid"},
					&cli.StringFlag{Name: "type", Aliases: []string{"t"}, Value: "limit", Usage: "limit"},
					&cli.StringFlag{Name: "price", Aliases: []string{"p"}, Value: "1.00", Usage: "价格小数点后随机"},
					&cli.StringFlag{Name: "qty", Aliases: []string{"q"}, Value: "10", Usage: "10"},
					&cli.IntFlag{Name: "n", DefaultText: "1", Value: 1, Usage: "循环插入订单个数"},
				},
				Action: func(ctx *cli.Context) error {
					if ctx.String("side") == "ask" {
						// haotrader.InsertAsk(rc, ctx.String("symbol"), ctx.String("type"), ctx.Int("n"), ctx.String("price"), ctx.String("qty"))
					} else {
						// haotrader.InsertBid(rc, ctx.String("symbol"), ctx.String("type"), ctx.Int("n"), ctx.String("price"), ctx.String("qty"))
					}
					return nil
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.Bool("deamon") {
				context, d, err := app.Deamon("haomatch.pid", "")
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

			matching.Run()
			return nil
		},
	}

	if err := appm.Run(os.Args); err != nil {
		app.Logger.Fatal(err)
	}
}
