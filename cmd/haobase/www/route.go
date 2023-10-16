package www

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/yzimhao/trading_engine/cmd/haobase/www/middle"
	"github.com/yzimhao/trading_engine/utils"
)

func Run() {
	demoBaseData()

	g := gin.New()
	router(g)
	g.Run(viper.GetString("haobase.http.host"))
}

func router(r *gin.Engine) {
	r.Use(utils.CorsMiddleware())

	api := r.Group("/api/v1/base")
	{

		api.Use(middle.CheckLogin())

		//todo 登陆验证
		api.GET("/assets/recharge", assets_recharge)
		api.GET("/assets", assets_balance)

		api.POST("/order/create", order_create)
		api.POST("/order/cancel", order_cancel)
		api.GET("/order/hisotry", nil)
		api.GET("/order/unfinished", nil)
	}

}
