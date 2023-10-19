package middle

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app/config"
)

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := ""
		if config.App.Main.Mode == config.ModeDemo {
			//自动为demo用户充值三种货币
			user_id = c.GetHeader("Token")
			if user_id == "" {
				user_id = c.Query("user_id")
			}

			if user_id != "" {
				pp := regexp.MustCompile(`^[a-z0-9]{4,10}$`)
				if !pp.MatchString(user_id) {
					utils.ResponseFailJson(c, "用户名不符合规则: ^[a-z0-9]{4,10}$")
					c.Abort()
					return
				}

				if assets.BalanceOfTotal(user_id, "usd").Equal(decimal.Zero) {
					assets.SysRecharge(user_id, "usd", "10000.00", "sys_recharge")
				}
				if assets.BalanceOfTotal(user_id, "jpy").Equal(decimal.Zero) {
					assets.SysRecharge(user_id, "jpy", "10000.00", "sys_recharge")
				}
				if assets.BalanceOfTotal(user_id, "eur").Equal(decimal.Zero) {
					assets.SysRecharge(user_id, "eur", "10000.00", "sys_recharge")
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
