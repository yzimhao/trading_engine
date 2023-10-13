package main

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/sevlyar/go-daemon"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"github.com/yzimhao/trading_engine/haotrader"
	"github.com/yzimhao/trading_engine/utils/app"
)

var (
	rc *redis.Client
)

func main() {
	appm := &cli.App{
		Name: "haotrader",
		// Version:   version,
		UsageText: "Issues: https://github.com/yzimhao/trading_engine/issues",
		Usage:     "交易撮合引擎",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Value: "./config.toml", Aliases: []string{"c"}},
			&cli.BoolFlag{Name: "deamon", Value: false, Aliases: []string{"d"}},
		},

		Before: func(cCtx *cli.Context) error {
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
					app.ConfigInit(ctx.String("config"))
					app.LogsInit("haotrader.run", false)
					rc := app.RedisInit()

					if ctx.String("side") == "ask" {
						haotrader.InsertAsk(rc, ctx.String("symbol"), ctx.String("type"), ctx.Int("n"), ctx.String("price"), ctx.String("qty"))
					} else {
						haotrader.InsertBid(rc, ctx.String("symbol"), ctx.String("type"), ctx.Int("n"), ctx.String("price"), ctx.String("qty"))
					}
					return nil
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			app.ConfigInit(ctx.String("config"))
			app.LogsInit("haotrader.run", ctx.Bool("deamon"))

			if viper.GetString("main.mode") != "prod" {
				logrus.Infof("当前运行在%s模式下，生产环境时main.mode请务必成prod", viper.GetString("main.mode"))
			}

			rc := app.RedisInit()
			ctext := context.Background()

			if ctx.Bool("deamon") {
				logrus.Info("开始守护进程")
				context, d, err := app.Deamon("haotrader.pid", "")
				if err != nil {
					logrus.Fatal("创建守护进程失败, err:", err.Error())
				}
				if d != nil {
					logrus.Printf("这是在父进程的标志")
					return nil
				}

				defer func(context *daemon.Context) {
					err := context.Release()
					if err != nil {
						logrus.Printf("释放失败:%s", err.Error())
					}
					logrus.Printf("释放成功!!!")
				}(context)
			}

			haotrader.Start(&ctext, rc)
			return nil
		},
	}

	if err := appm.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
