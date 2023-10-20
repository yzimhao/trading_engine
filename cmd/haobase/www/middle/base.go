package middle

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
	"github.com/yzimhao/trading_engine/utils/app/config"
)

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := ""
		if config.App.Main.Mode == config.ModeDemo {
			user_id = c.GetHeader("Token")
			if user_id == "" {
				user_id = c.Query("user_id")
			}

			app.Logger.Infof("[%s]: %s %s", c.ClientIP(), c.Request.Method, c.Request.RequestURI)

			if user_id != "" {
				pp := regexp.MustCompile(`^[a-z0-9]{4,10}$`)
				if !pp.MatchString(user_id) {
					utils.ResponseFailJson(c, "用户名不符合规则: ^[a-z0-9]{4,10}$")
					c.Abort()
					return
				}
			}
		}

		if user_id == "" {
			utils.ResponseFailJson(c, "需要登录")
			c.Abort()
			return
		}

		c.Set("user_id", user_id)

		// before request
		c.Next()

	}
}
