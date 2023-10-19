package main

import (
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/gookit/goutil/arrutil"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"github.com/yzimhao/trading_engine/utils/app"

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
			app.ConfigInit(ctx.String("config"), ctx.Bool("deamon"))
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
				context, d, err := app.Deamon("run.pid", "")
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

			startWeb(viper.GetString("example.host"))
			return nil
		},
	}

	if err := appm.Run(os.Args); err != nil {
		app.Logger.Fatal(err)
	}
}

func startWeb(host string) {
	web = gin.New()
	web.SetFuncMap(template.FuncMap{
		"unsafe": func(str string) template.HTML {
			return template.HTML(str)
		},
	})

	web.LoadHTMLGlob("./*.html")
	web.StaticFS("/statics", http.Dir("./statics"))

	web.GET("/:symbol", func(c *gin.Context) {
		support := []string{"usdjpy", "eurusd"}
		symbol := strings.ToLower(c.Param("symbol"))

		if !arrutil.Contains(support, symbol) {
			c.Redirect(301, "/")
			return
		}

		c.HTML(200, "demo.html", gin.H{
			"haobase_host":  viper.GetString("api.haobase_host"),
			"haoquote_host": viper.GetString("api.haoquote_host"),
			"ws_host":       viper.GetString("api.haoquote_ws_host"),
			"symbol":        symbol,
		})
	})

	web.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/usdjpy")
	})

	web.Run(host)
}
