package middlewares

import "github.com/gin-gonic/gin"

type Middleware struct{}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userId", "123")
		c.Next()
	}
}
