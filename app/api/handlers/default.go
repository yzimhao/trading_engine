package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yzimhao/trading_engine/v2/app/api/handlers/controllers"
	_ "github.com/yzimhao/trading_engine/v2/app/docs"
)

type RoutesContext struct {
	engine *gin.Engine
	//middleware
	//controllers
	userAssetsController *controllers.UserAssetsController
}

func NewRoutesContext(
	engine *gin.Engine,
	userAssetsController *controllers.UserAssetsController,
) *RoutesContext {
	r := &RoutesContext{
		engine:               engine,
		userAssetsController: userAssetsController,
	}

	r.registerRoutes()
	return r
}

func (ctx *RoutesContext) registerRoutes() {
	ctx.engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	ctx.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiGroup := ctx.engine.Group("api")
	v1Group := apiGroup.Group("v1")

	//for test
	wallet := v1Group.Group("wallet")
	wallet.POST("/assets/despoit", ctx.userAssetsController.Despoit)
	wallet.POST("/assets/withdraw", ctx.userAssetsController.Withdraw)
	wallet.GET("/assets/:symbol", ctx.userAssetsController.Query)
	wallet.GET("/assets/:symbol/history", ctx.userAssetsController.QueryAssetHistory)
	wallet.POST("/transfer/:symbol", ctx.userAssetsController.Transfer)

}
