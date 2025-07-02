package provider

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func (r *Router) ResponseOk(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"ok":   true,
		"data": data,
	})
}

// todo 携带错误码
func (r *Router) ResponseError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"ok":  false,
		"msg": err.Error(),
	})
}
