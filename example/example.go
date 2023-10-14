package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gookit/goutil/arrutil"
	"github.com/sevlyar/go-daemon"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"github.com/yzimhao/trading_engine/utils/app"
	"github.com/zehuamama/balancer/proxy"

	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
)

var (
	web *gin.Engine
)

func main() {

	appm := &cli.App{
		Name:      "example",
		UsageText: "Issues: https://github.com/yzimhao/trading_engine/issues",
		Usage:     "交易系统示例程序",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Value: "./config.toml", Aliases: []string{"c"}},
			&cli.BoolFlag{Name: "deamon", Value: false, Aliases: []string{"d"}},
		},

		Before: func(ctx *cli.Context) error {
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
			app.ConfigInit(ctx.String("config"))
			app.LogsInit("run.log", ctx.Bool("deamon"))

			if viper.GetString("main.mode") != "prod" {
				logrus.Infof("当前运行在%s模式下，生产环境时main.mode请务必成prod", viper.GetString("main.mode"))
			}

			if ctx.Bool("deamon") {
				context, d, err := app.Deamon("run.pid", "")
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

			startWeb(viper.GetString("example.host"))
			return nil
		},
	}

	if err := appm.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func startWeb(host string) {
	web = gin.New()
	web.LoadHTMLGlob("./*.html")
	web.StaticFS("/statics", http.Dir("./statics"))

	//代理交易系统的后端接口，实际应用中可以用nginx直接代理
	web.Any("/api/v1/base/*any", func(ctx *gin.Context) {
		httpProxy, err := proxy.NewHTTPProxy([]string{viper.GetString("api.haobase_host")}, "round-robin")
		if err != nil {
			log.Fatalf("create proxy error: %s", err)
		}
		httpProxy.ServeHTTP(ctx.Writer, ctx.Request)
	})
	web.Any("/api/v1/quote/*any", func(ctx *gin.Context) {
		httpProxy, err := proxy.NewHTTPProxy([]string{viper.GetString("api.haoquote_host")}, "round-robin")
		if err != nil {
			log.Fatalf("create proxy error: %s", err)
		}
		httpProxy.ServeHTTP(ctx.Writer, ctx.Request)
	})

	web.GET("/:symbol", func(c *gin.Context) {
		support := []string{"usdjpy", "eurusd"}
		symbol := strings.ToLower(c.Param("symbol"))

		if !arrutil.Contains(support, symbol) {
			c.Redirect(301, "/")
			return
		}

		c.HTML(200, "demo.html", gin.H{
			"symbol": symbol,
		})
	})

	web.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/usdjpy")
	})

	web.Run(host)
}
