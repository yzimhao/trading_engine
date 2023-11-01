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

type tradlogSearch struct {
	Symbol string `json:"symbol"`
}

func TradeHistory(ctx *gin.Context) {
	db := app.Database().NewSession()
	defer db.Close()

	if ctx.Request.Method == "GET" {
		page := utils.S2Int(ctx.Query("page"))
		limit := utils.S2Int(ctx.Query("limit"))
		searchParams := ctx.Query("searchParams")

		var search tradlogSearch
		json.Unmarshal([]byte(searchParams), &search)

		if page <= 0 {
			page = 1
		}
		if limit <= 0 {
			limit = 10
		}
		offset := (page - 1) * limit

		data := []orders.TradeLog{}

		if search.Symbol == "" {
			for _, item := range base.NewTSymbols().All() {
				tb := orders.GetTradelogTableName(item.Symbol)
				if dbtables.Exist(db, tb) {
					search.Symbol = item.Symbol
					break
				}
			}
		}

		tablename := orders.GetTradelogTableName(search.Symbol)
		q := db.Table(tablename)
		cond := q.Conds()
		err := q.OrderBy("id desc").Limit(limit, offset).Find(&data)
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
			ctx.HTML(200, "user_trade_history", gin.H{
				"searchParams": searchParams,
				"all_symbols":  base.NewTSymbols().All(),
			})
		}
		return
	}
}
