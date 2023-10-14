package middle

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := ""
		if app.RunMode == app.ModeDemo {
			user_id = c.GetHeader("UserId")
		}

		if user_id == "" {
			utils.ResponseFailJson(c, "需要登录")
			return
		}

		c.Set("user_id", user_id)

		// before request
		c.Next()

	}
}
