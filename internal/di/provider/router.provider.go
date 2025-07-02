package provider

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/internal/types"
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
		"code": types.SuccessCode,
		"data": data,
	})
}

// 携带错误码
func (r *Router) ResponseError(c *gin.Context, code types.ErrorCode) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  types.GetErrorMsg(code),
	})
}

// 从jwt认证中获取用户ID
func (r *Router) ParseUserID(c *gin.Context) string {
	claims := jwt.ExtractClaims(c)
	return claims["userId"].(string)
}
