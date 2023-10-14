package middle

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := ""
		if app.RunMode == app.ModeDemo {
			//自动为demo用户充值三种货币
			user_id = c.GetHeader("Token")
			if user_id != "" {
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
