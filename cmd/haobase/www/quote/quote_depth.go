package quote

import (
	"strings"

	"github.com/gin-gonic/gin"
	haoquote "github.com/yzimhao/trading_engine/cmd/haoquote/quote"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func QuoteDepth(ctx *gin.Context) {
	limit := utils.S2Int(ctx.Query("limit"))
	symbol := strings.ToLower(ctx.Query("symbol"))

	rdc := app.RedisPool().Get()
	defer rdc.Close()

	data, err := haoquote.GetDepthData(symbol)
	if err != nil {
		utils.ResponseFailJson(ctx, "invalid symbol")
		return
	}

	if limit <= 0 || limit > 50 {
		limit = 10
	}

	utils.ResponseOkJson(ctx, gin.H{
		"asks":    haoquote.LimitSize(data.Asks, limit),
		"bids":    haoquote.LimitSize(data.Bids, limit),
		"asksize": data.AsksSize,
		"bidsize": data.BidsSize,
	})
}

func QuoteLatestPrice(ctx *gin.Context) {
	symbol := strings.ToLower(ctx.Query("symbol"))

	rdc := app.RedisPool().Get()
	defer rdc.Close()

	data, err := haoquote.GetDepthData(symbol)
	if err != nil {
		utils.ResponseFailJson(ctx, "invalid symbol")
		return
	}

	utils.ResponseOkJson(ctx, gin.H{
		symbol: data.Price,
	})
}
