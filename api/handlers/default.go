package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/api/handlers/controllers"
)

type RoutesContext struct {
	engine *gin.Engine
	//middleware
	//controllers
	userAssetsController *controllers.UserAssetsController
	orderController      *controllers.OrderController
}

func NewRoutesHandler(engine *gin.Engine) {
	r := &RoutesContext{
		engine: engine,
	}

	r.registerRoutes()
}

func (ctx *RoutesContext) registerRoutes() {
	ctx.engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	apiGroup := ctx.engine.Group("api")
	v1Group := apiGroup.Group("v1")

	v1Group.Group("user").
		POST("", ctx.userAssetsController.Create).
		PATCH("", ctx.userAssetsController.Update).
		DELETE("", ctx.userAssetsController.Delete)

}
