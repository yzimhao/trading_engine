package order

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haobase/orders"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func Unfinished(ctx *gin.Context) {
	user_id := ctx.MustGet("user_id").(string)
	limit := utils.S2Int(ctx.Query("limit"))
	symbol := ctx.Query("symbol")

	cfg, err := base.NewTSymbols().Get(symbol)
	if err != nil {
		utils.ResponseFailJson(ctx, "不存在的交易对")
		return
	}

	db := app.Database().NewSession()
	defer db.Close()

	var rows []orders.Order
	query := db.Table(new(orders.UnfinishedOrder))

	if symbol != "" {
		query = query.And("symbol=?", symbol)
	}
	query.Where("user_id=?", user_id).OrderBy("create_time desc").Limit(limit).Find(&rows)

	for i, v := range rows {
		rows[i] = v.FormatDecimal(cfg.PricePrecision, cfg.QtyPrecision)
	}

	utils.ResponseOkJson(ctx, rows)
}
