package order

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haobase/orders"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func Unfinished(ctx *gin.Context) {
	user_id := ctx.MustGet("user_id").(string)

	symbol := ctx.Query("symbol")

	db := app.Database().NewSession()
	defer db.Close()

	var rows []orders.UnfinishedOrder
	query := db.Table(new(orders.UnfinishedOrder))

	if symbol != "" {
		query = query.And("symbol=?", symbol)
	}
	query.Where("user_id=?", user_id).OrderBy("create_time desc").Find(&rows)

	utils.ResponseOkJson(ctx, rows)
}
