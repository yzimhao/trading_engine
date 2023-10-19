package www

import (
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

		api.Use(middle.CheckLogin())

		api.GET("/assets", assets_balance)

		api.POST("/order/create", order.Create)
		api.POST("/order/cancel", order.Cancel)
		api.GET("/order/hisotry", order.History)
		api.GET("/order/unfinished", order.Unfinished)
	}

}
