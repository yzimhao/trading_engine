package www

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/utils"
)

func assets_recharge(ctx *gin.Context) {
	// assets.SysRecharge("user1", "usd", "10000.00", "C001")
	// assets.SysRecharge("user1", "jpy", "10000.00", "C002")

	// assets.SysRecharge("user2", "usd", "10000.00", "C001")
	// assets.SysRecharge("user2", "jpy", "10000.00", "C002")
}

func assets_balance(ctx *gin.Context) {
	ss := ctx.Query("symbols")

	symbols := make([]string, 0)
	if ss != "" {
		symbols = strings.Split(ss, ",")
	}
	rows := assets.UserAssets("user1", symbols)
	utils.ResponseOkJson(ctx, rows)
}
