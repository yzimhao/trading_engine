package www

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func qutoe_depth(ctx *gin.Context) {
	limit := utils.S2Int(ctx.Query("limit"))
	symbol := strings.ToLower(ctx.Query("symbol"))

	rdc := app.RedisPool().Get()
	defer rdc.Close()

	data, err := get_depth_data(symbol)
	if err != nil {
		utils.ResponseFailJson(ctx, "invalid symbol")
		return
	}

	if limit <= 0 || limit > 100 {
		limit = 10
	}

	utils.ResponseOkJson(ctx, gin.H{
		"asks": data.Asks, // symbols_depth.limit("asks", symbol, limit),
		"bids": data.Bids, //symbols_depth.limit("bids", symbol, limit),
	})
}
