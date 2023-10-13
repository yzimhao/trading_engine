package www

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Run() {
	g := gin.New()
	router(g)
	g.Run(viper.GetString("haobase.http.host"))
}

func router(r *gin.Engine) {
	api := r.Group("/api/v1/base")
	{
		//todo 登陆验证
		api.GET("/assets/recharge", assets_recharge)
		api.GET("/assets", assets_balance)

		api.POST("/order/create", order_create)
		api.POST("/order/cancel", order_cancel)
		api.GET("/order/hisotry", nil)
		api.GET("/order/unfinished", nil)
	}

}
