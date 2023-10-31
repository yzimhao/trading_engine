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
				tb := orders.GetOrderTableName(item.Symbol)
				if dbtables.Exist(db, tb) {
					search.Symbol = item.Symbol
					break
				}
			}
		}

		tablename := orders.GetOrderTableName(search.Symbol)
		q := db.Table(tablename)
		q = q.Where("symbol = ?", search.Symbol)
		cond := q.Conds()
		err := q.OrderBy("id desc").Limit(limit, offset).Find(&data)
		if err != nil {
			render(ctx, 1, err.Error(), 0, "")
			return
		}

		total, _ := q.Table(tablename).And(cond).Count()
		if ctx.Query("api") == "1" {
			render(ctx, 0, "", int(total), data)
		} else {
			ctx.HTML(200, "user_order_history", gin.H{
				"searchParams": searchParams,
				"all_symbols":  base.NewTSymbols().All(),
			})
		}
		return
	}
}

type orderSearch struct {
	Symbol string `json:"symbol"`
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
			q = q.Where("symbol like ?", "%"+search.Symbol+"%")
		}

		cond := q.Conds()
		err := q.OrderBy("id desc").Limit(limit, offset).Find(&data)
		if err != nil {
			render(ctx, 1, err.Error(), 0, "")
			return
		}

		total, _ := db.Table(new(orders.UnfinishedOrder)).And(cond).Count()

		for i, v := range data {
			cfg, _ := base.NewTSymbols().Get(v.Symbol)
			data[i].Price = utils.FormatDecimal(v.Price, cfg.PricePrecision)
			data[i].Quantity = utils.FormatDecimal(v.Quantity, cfg.QtyPrecision)
			data[i].FinishedQty = utils.FormatDecimal(v.FinishedQty, cfg.QtyPrecision)
			data[i].FinishedAmount = utils.FormatDecimal(v.FinishedAmount, cfg.PricePrecision)
		}

		if ctx.Query("api") == "1" {
			render(ctx, 0, "", int(total), data)
		} else {
			ctx.HTML(200, "user_unfinished", gin.H{
				"searchParams": searchParams,
			})
		}
		return
	}
}
