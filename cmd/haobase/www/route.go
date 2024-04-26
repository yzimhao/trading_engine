package www

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haobase/message"
	"github.com/yzimhao/trading_engine/cmd/haobase/message/ws"
	"github.com/yzimhao/trading_engine/cmd/haobase/www/internal_api"
	"github.com/yzimhao/trading_engine/cmd/haobase/www/middle"
	"github.com/yzimhao/trading_engine/cmd/haobase/www/order"
	"github.com/yzimhao/trading_engine/cmd/haobase/www/quote"
	"github.com/yzimhao/trading_engine/config"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func Run() {
	if config.App.Main.Mode != app.ModeProd.String() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	g := gin.New()
	router(g)
	app.Logger.Infof("listen: %s", config.App.Haobase.Listen)
	g.Run(config.App.Haobase.Listen)
}

func router(r *gin.Engine) {
	r.Use(utils.CorsMiddleware())

	router_internal(r)
	router_base(r)
	router_quote(r)
	router_wss(r)
}

func router_internal(r *gin.Engine) {
	internal := r.Group("/api/v1/internal")
	{ //内部通信接口
		//todo 加入ip限制
		internal.Use(internal_api.Authentication())

		internal.POST("/settoken", internal_api.SetToken)
		//内部充值
		internal.POST("/deposit", internal_api.Deposit)
		internal.POST("/withdraw", internal_api.Withdraw)
	}
}

func router_base(r *gin.Engine) {
	base_api := r.Group("/api/v1/base")
	{
		base_api.GET("/ping", func(ctx *gin.Context) {
			utils.ResponseOkJson(ctx, gin.H{})
		})
		base_api.GET("/time", func(ctx *gin.Context) {
			utils.ResponseOkJson(ctx, gin.H{
				"server_time": time.Now().Unix(),
			})
		})
		base_api.GET("/version", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"version": app.Version,
				"build":   app.Build,
			})
		})

		//全部交易品类
		base_api.GET("/trading/varieties", trading_varieties)
		//指定交易品类
		base_api.GET("/varieties/config", varieties_config)

		base_api.Use(middle.CheckLogin())
		{
			if config.App.Main.Mode == app.ModeDemo.String() {
				base_api.GET("/assets/recharge_for_demo", func(ctx *gin.Context) {
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
			base_api.GET("/assets", assets_balance)
			base_api.POST("/order/create", order.Create)
			base_api.POST("/order/cancel", order.Cancel)
			base_api.GET("/order/history", order.History)
			base_api.GET("/order/unfinished", order.Unfinished)
		}
	}
}

func router_wss(r *gin.Engine) {
	ws.NewHub()
	message.Subscribe()

	r.GET("/ws", func(ctx *gin.Context) {
		ws.M.ServeWs(ctx)
	})
}

func router_quote(r *gin.Engine) {

	quote_api := r.Group("/api/v1/quote")
	quote_api.Use(utils.CorsMiddleware())
	{
		quote_api.GET("/depth", quote.QuoteDepth)
		quote_api.GET("/price", quote.QuoteLatestPrice)
		quote_api.GET("/trans/record", quote.TransRecord)
		quote_api.GET("/kline", quote.KLine)
	}
}
