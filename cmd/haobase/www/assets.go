package www

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/utils"
)

func assets_balance(ctx *gin.Context) {
	user_id := ctx.MustGet("user_id").(string)
	ss := ctx.Query("symbols")

	symbols := make([]string, 0)
	if ss != "" {
		symbols = strings.Split(ss, ",")
	}
	rows := assets.UserAssets(user_id, symbols)
	utils.ResponseOkJson(ctx, rows)
}
