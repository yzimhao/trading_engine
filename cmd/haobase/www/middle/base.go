package middle

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/utils/app"
)

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", "")
		if app.RunMode == app.ModeDemo {
			user_id := c.GetHeader("user_id")
			c.Set("user_id", user_id)
		}

		// before request
		c.Next()

	}
}
