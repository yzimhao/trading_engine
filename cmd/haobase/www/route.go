package www

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haobase/www/middle"
	"github.com/yzimhao/trading_engine/cmd/haobase/www/order"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
	"github.com/yzimhao/trading_engine/utils/app/config"
)

func Run() {
	if !config.App.Haobase.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	g := gin.New()
	router(g)
	app.Logger.Infof("listen: %s", config.App.Haobase.Listen)
	g.Run(config.App.Haobase.Listen)
}

func router(r *gin.Engine) {
	r.Use(utils.CorsMiddleware())

	api := r.Group("/api/v1/base")
	{
		api.GET("/ping", func(ctx *gin.Context) {
			utils.ResponseOkJson(ctx, gin.H{})
		})
		api.GET("/time", func(ctx *gin.Context) {
			utils.ResponseOkJson(ctx, gin.H{
				"server_time": time.Now().Unix(),
			})
		})
		//全部交易品类
		api.GET("/trading/varieties", trading_varieties)
		//指定交易品类
		api.GET("/varieties/config", varieties_config)

		api.Use(middle.CheckLogin())
		{
			if config.App.Main.Mode == config.ModeDemo {
				api.GET("/assets/recharge_for_demo", recharge_for_demo)
			}

			api.GET("/assets", assets_balance)
			api.POST("/order/create", order.Create)
			api.POST("/order/cancel", order.Cancel)
			api.GET("/order/history", order.History)
			api.GET("/order/unfinished", order.Unfinished)
		}
	}

}
