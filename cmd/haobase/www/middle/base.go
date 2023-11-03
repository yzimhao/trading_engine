package middle

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gookit/goutil/arrutil"
	"github.com/yzimhao/trading_engine/cmd/haobase/www/internal_api"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
	"github.com/yzimhao/trading_engine/utils/app/config"
)

var (
	sysUsers = []string{"root", "fee"}
)

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := ""
		token := c.GetHeader("Token")
		if config.App.Main.Mode == config.ModeDemo {
			user_id = token
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

				if arrutil.Contains(sysUsers, user_id) {
					utils.ResponseFailJson(c, "禁止登陆")
					c.Abort()
					return
				}
			}
		} else {
			//从redis的token中获取登陆用户ID
			user_id = internal_api.GetUserIdFromToken(token)
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
