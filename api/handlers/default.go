package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoutesContext struct {
	engine *gin.Engine
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
}
