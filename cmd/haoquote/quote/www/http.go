package www

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/yzimhao/trading_engine/cmd/haobase/message"
	"github.com/yzimhao/trading_engine/cmd/haobase/message/ws"
	_ "github.com/yzimhao/trading_engine/docs/api" // main 文件中导入 docs 包
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
	"github.com/yzimhao/trading_engine/utils/app/config"
)

var ()

func Run() {
	publish_depth()
	http_start()
}

func http_start() {
	if config.App.Haoquote.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	logrus.Infof("HTTP服务: %s", config.App.Haoquote.Listen)
	g := gin.New()
	router(g)
	g.Run(config.App.Haoquote.Listen)
}

func router(router *gin.Engine) {
	//websokect服务放在这个quote里
	ws.NewHub()
	message.Subscribe()

	router.GET("/quote/ws", func(ctx *gin.Context) {
		ws.M.ServeWs(ctx)
	})

	api := router.Group("/api/v1/quote")
	{
		api.Use(utils.CorsMiddleware())
		api.GET("/depth", qutoe_depth)
		api.GET("/trans/record", trans_record)
		api.GET("/kline", kline)
		api.GET("/system", system_info)
	}
}

func system_info(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"version": app.Version,
		"build":   app.Build,
	})
}
