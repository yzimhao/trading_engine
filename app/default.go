package app

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yzimhao/trading_engine/v2/app/api/handlers/controllers"
	_ "github.com/yzimhao/trading_engine/v2/app/docs"
	"github.com/yzimhao/trading_engine/v2/app/middlewares"
	"github.com/yzimhao/trading_engine/v2/app/webws"
	"go.uber.org/fx"
)

type Routes struct {
	engine *gin.Engine
	//middleware
	//controllers
	baseController       *controllers.BaseController
	userAssetsController *controllers.UserAssetsController
	orderController      *controllers.OrderController
	marketController     *controllers.MarketController
	authMiddleware       *middlewares.AuthMiddleware
	userController       *controllers.UserController
	wsController         *webws.WsManager
}

type inContext struct {
	fx.In
	Engine               *gin.Engine
	BaseController       *controllers.BaseController
	UserAssetsController *controllers.UserAssetsController
	OrderController      *controllers.OrderController
	MarketController     *controllers.MarketController
	AuthMiddleware       *middlewares.AuthMiddleware
	UserController       *controllers.UserController
	WsController         *webws.WsManager
}

func NewRoutes(in inContext) *Routes {
	r := &Routes{
		engine:               in.Engine,
		baseController:       in.BaseController,
		userAssetsController: in.UserAssetsController,
		orderController:      in.OrderController,
		marketController:     in.MarketController,
		authMiddleware:       in.AuthMiddleware,
		userController:       in.UserController,
		wsController:         in.WsController,
	}

	r.registerRoutes()
	return r
}

func (ctx *Routes) registerRoutes() {

	ctx.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//ws
	ctx.engine.GET("/ws", func(c *gin.Context) {
		ctx.wsController.Listen(c.Writer, c.Request, c.Request.Header)
	})

	apiGroup := ctx.engine.Group("api")
	v1Group := apiGroup.Group("v1")
	// v1Group.GET("/ping", ctx.baseController.Ping)
	// v1Group.GET("/time", ctx.baseController.Time)

	base := v1Group.Group("base")
	base.GET("/exchange_info", ctx.baseController.ExchangeInfo)

	user := v1Group.Group("user")
	user.POST("/login", ctx.authMiddleware.Jwt().LoginHandler)
	user.POST("/refresh_token", ctx.authMiddleware.Jwt().RefreshHandler)

	asset := user.Group("asset")
	asset.Use(ctx.authMiddleware.Auth())
	asset.POST("/despoit", ctx.userAssetsController.Despoit)
	asset.POST("/withdraw", ctx.userAssetsController.Withdraw)
	asset.GET("/query", ctx.userAssetsController.Query)
	asset.GET("/:symbol/history", ctx.userAssetsController.QueryAssetHistory)
	asset.POST("/transfer/:symbol", ctx.userAssetsController.Transfer)

	order := v1Group.Group("order")
	order.Use(ctx.authMiddleware.Auth())
	order.POST("/create", ctx.orderController.Create)
	order.GET("/history", ctx.orderController.HistoryList)
	order.GET("/unfinished", ctx.orderController.UnfinishedList)
	order.GET("/trade/history", ctx.orderController.TradeHistoryList)

	market := v1Group.Group("market")
	market.GET("/depth", ctx.marketController.Depth)
	market.GET("/trades", ctx.marketController.Trades)
	market.GET("/klines", ctx.marketController.Klines)
}
