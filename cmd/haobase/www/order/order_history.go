package order

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haobase/orders"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func History(ctx *gin.Context) {
	user_id := ctx.MustGet("user_id").(string)
	limit := utils.S2Int(ctx.Query("limit"))
	symbol := ctx.Query("symbol")

	cfg, err := base.NewTSymbols().Get(symbol)
	if err != nil {
		utils.ResponseFailJson(ctx, err.Error())
		return
	}

	db := app.Database().NewSession()
	defer db.Close()

	rows := make([]orders.Order, 0)
	table := &orders.Order{Symbol: symbol}
	query := db.Table(table)

	query.Where("user_id=? and status !=?", user_id, orders.OrderStatusNew).OrderBy("create_time desc").Limit(limit).Find(&rows)

	for i, v := range rows {
		rows[i] = v.FormatDecimal(cfg.PricePrecision, cfg.QtyPrecision)
	}

	utils.ResponseOkJson(ctx, rows)
}
