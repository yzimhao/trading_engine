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
		"asks": limitSize(data.Asks, limit),
		"bids": limitSize(data.Bids, limit),
	})
}

func limitSize(arr [][2]string, n int) [][2]string {
	a := len(arr)
	if n >= a {
		n = a
	}
	if n <= 0 {
		n = 0
	}
	return arr[0:n]
}
