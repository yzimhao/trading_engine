package admin

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haobase/orders"
	"github.com/yzimhao/trading_engine/types/dbtables"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

type orderSearch struct {
	Symbol    string `json:"symbol"`
	Status    string `json:"status"`
	UserId    string `json:"user_id"`
	OrderType string `json:"order_type"`
	OrderId   string `json:"order_id"`
}

func UserOrderHistory(ctx *gin.Context) {
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

		if search.Symbol == "" {
			for _, item := range base.NewTSymbols().All() {
				bean := orders.Order{Symbol: item.Symbol}
				if dbtables.Exist(db, &bean) {
					search.Symbol = item.Symbol
					break
				}
			}
		}

		tablename := &orders.Order{Symbol: search.Symbol}
		q := db.Table(tablename)
		q = q.Where("symbol = ? and status>?", search.Symbol, orders.OrderStatusNew)

		if search.UserId != "" {
			q = q.Where("user_id = ?", search.UserId)
		}

		if search.OrderType != "" {
			q = q.Where("order_type = ?", search.OrderType)
		}

		if search.OrderId != "" {
			q = q.Where("order_id = ?", search.OrderId)
		}

		cond := q.Conds()
		err := q.OrderBy("create_time desc").Limit(limit, offset).Find(&data)
		if err != nil {
			render(ctx, 1, err.Error(), 0, "")
			return
		}

		total, _ := q.Table(tablename).And(cond).Count()

		cfg, _ := base.NewTSymbols().Get(search.Symbol)
		for i, _ := range data {
			data[i].FormatDecimal(cfg.PricePrecision, cfg.QtyPrecision)
		}

		if ctx.Query("api") == "1" {
			render(ctx, 0, "", int(total), data)
		} else {
			ctx.HTML(200, "user_order_history", gin.H{
				"search":      search,
				"all_symbols": base.NewTSymbols().All(),
			})
		}
		return
	}
}
