package www

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/utils"
)

// 全部交易品类
func trading_varieties(ctx *gin.Context) {
	utils.ResponseOkJson(ctx, base.NewTSymbols().All())
}

func varieties_config(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	data, _ := base.NewTSymbols().Get(symbol)
	utils.ResponseOkJson(ctx, data)
}
