package admin

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haobase/orders"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func CancelUserOrder(ctx *gin.Context) {
	order := ctx.PostForm("order_ids")
	ids := strings.Split(order, ",")
	for _, order_id := range ids {
		if err := orders.SubmitOrderCancel(order_id, trading_core.CancelTypeBySystem); err != nil {
			logrus.Errorf("CancelUserOrder %s err: %s", order_id, err.Error())
		}
	}
	utils.ResponseOkJson(ctx, "")
}

func UserOrderUnfinished(ctx *gin.Context) {
	db := app.Database().NewSession()
	defer db.Close()

	if ctx.Request.Method == "GET" {
		page := utils.S2Int(ctx.Query("page"))
		limit := utils.S2Int(ctx.Query("limit"))
		searchParams := ctx.Query("searchParams")

		var search orderSearch
		json.Unmarshal([]byte(searchParams), &search)

		if page <= 0 {
			page = 1
		}
		if limit <= 0 {
			limit = 10
		}
		offset := (page - 1) * limit

		data := []orders.Order{}

		q := db.Table(new(orders.UnfinishedOrder))

		if search.Symbol != "" {
			q = q.Where("symbol = ?", search.Symbol)
		}
		if search.UserId != "" {
			q = q.Where("user_id = ?", search.UserId)
		}
		if search.OrderId != "" {
			q = q.Where("order_id = ?", search.OrderId)
		}

		cond := q.Conds()
		err := q.OrderBy("price desc,create_time asc").Limit(limit, offset).Find(&data)
		if err != nil {
			render(ctx, 1, err.Error(), 0, "")
			return
		}

		total, _ := db.Table(new(orders.UnfinishedOrder)).And(cond).Count()

		for i, v := range data {
			cfg, _ := base.NewTradeSymbol().Get(v.Symbol)
			data[i].FormatDecimal(cfg.PricePrecision, cfg.QtyPrecision)
		}

		if ctx.Query("api") == "1" {
			render(ctx, 0, "", int(total), data)
		} else {
			ctx.HTML(200, "user_unfinished", gin.H{
				"search":      search,
				"all_symbols": base.NewTradeSymbol().All(),
			})
		}
		return
	}
}
