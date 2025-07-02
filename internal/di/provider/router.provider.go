package provider

import "github.com/gin-gonic/gin"

type Router struct {
	APIv1 *gin.RouterGroup
	*gin.Engine
}

func NewRouter(engine *gin.Engine) *Router {
	return &Router{
		Engine: engine,
		APIv1:  engine.Group("/api/v1"),
	}
}
