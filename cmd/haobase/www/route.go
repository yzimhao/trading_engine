package www

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haobase/www/internal_api"
	"github.com/yzimhao/trading_engine/cmd/haobase/www/middle"
	"github.com/yzimhao/trading_engine/cmd/haobase/www/order"
	"github.com/yzimhao/trading_engine/config"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
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

	internal := r.Group("/api/v1/internal")
	{ //内部通信接口
		//todo 加入ip限制
		internal.Use(internal_api.Authentication())

		internal.POST("/settoken", internal_api.SetToken)
		//内部充值
		internal.POST("/deposit", internal_api.Deposit)
		internal.POST("/withdraw", internal_api.Withdraw)
	}

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
			if config.App.Main.Mode == app.ModeDemo.String() {
				api.GET("/assets/recharge_for_demo", func(ctx *gin.Context) {
					user_id := ctx.MustGet("user_id").(string)
					//自动为demo用户充值
					default_amount := "10000.00"
					all := base.NewSymbols().All()
					for _, item := range all {
						if assets.BalanceOfTotal(user_id, item.Symbol).Equal(decimal.Zero) {
							assets.SysDeposit(user_id, item.Symbol, default_amount, "sys.give:"+user_id)
						}
					}
					utils.ResponseOkJson(ctx, "")
				})
			}
			api.GET("/assets", assets_balance)
			api.POST("/order/create", order.Create)
			api.POST("/order/cancel", order.Cancel)
			api.GET("/order/history", order.History)
			api.GET("/order/unfinished", order.Unfinished)
		}
	}

}
