package admin

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haobase/base/varieties"
	"github.com/yzimhao/trading_engine/types"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func VarietiesAdd(ctx *gin.Context) {
	id := utils.S2Int(ctx.Query("id"))

	db := app.Database().NewSession()
	defer db.Close()

	if ctx.Request.Method == "GET" {
		data := varieties.Varieties{Id: id}

		if id > 0 {
			db.Table(new(varieties.Varieties)).Get(&data)
		}

		ctx.HTML(200, "varieties_add", gin.H{
			"data": data,
		})
	} else {
		data := varieties.Varieties{
			Id: id,
			// Symbol:        ctx.PostForm("symbol"),
			Name:          ctx.PostForm("name"),
			MinPrecision:  utils.S2Int(ctx.PostForm("min_precision")),
			ShowPrecision: utils.S2Int(ctx.PostForm("show_precision")),
			Sort:          utils.S2Int64(ctx.PostForm("sort")),
			Status:        types.ParseStatusString(ctx.PostForm("status")),
		}

		var err error
		if id > 0 {
			_, err = db.Table(new(varieties.Varieties)).Where("id=?", id).Cols("name,min_precision,show_precision,sort,status").Update(&data)
		} else {
			data.Symbol = strings.Trim(ctx.PostForm("symbol"), " ")
			_, err = db.Table(new(varieties.Varieties)).Insert(&data)
		}

		if err != nil {
			utils.ResponseFailJson(ctx, err.Error())
			return
		}
		utils.ResponseOkJson(ctx, "")
	}
}

type varietiesSearch struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func VarietiesList(ctx *gin.Context) {
	db := app.Database().NewSession()
	defer db.Close()

	if ctx.Request.Method == "GET" {
		page := utils.S2Int(ctx.Query("page"))
		limit := utils.S2Int(ctx.Query("limit"))
		searchParams := ctx.Query("searchParams")

		var search varietiesSearch
		json.Unmarshal([]byte(searchParams), &search)

		if page <= 0 {
			page = 1
		}
		if limit <= 0 {
			limit = 10
		}
		offset := (page - 1) * limit

		data := []varieties.Varieties{}

		q := db.Table(new(varieties.Varieties))

		if search.Symbol != "" {
			q = q.Where("symbol like ?", "%"+search.Symbol+"%")
		}
		if search.Name != "" {
			q = q.Where("name like ?", "%"+search.Name+"%")
		}
		if search.Status != "" {
			q = q.Where("status = ?", search.Status)
		}

		cond := q.Conds()
		err := q.OrderBy("sort asc, id desc").Limit(limit, offset).Find(&data)
		if err != nil {
			render(ctx, 1, err.Error(), 0, "")
			return
		}

		total, _ := db.Table(new(varieties.Varieties)).And(cond).Count()
		if ctx.Query("api") == "1" {
			render(ctx, 0, "", int(total), data)
		} else {
			ctx.HTML(200, "varieties_list", gin.H{
				"searchParams": searchParams,
			})
		}
		return
	}
}

func TradingVarietiesList(ctx *gin.Context) {
	db := app.Database().NewSession()
	defer db.Close()

	if ctx.Request.Method == "GET" {
		page := utils.S2Int(ctx.Query("page"))
		limit := utils.S2Int(ctx.Query("limit"))
		searchParams := ctx.Query("searchParams")

		var search varietiesSearch
		json.Unmarshal([]byte(searchParams), &search)

		if page <= 0 {
			page = 1
		}
		if limit <= 0 {
			limit = 10
		}
		offset := (page - 1) * limit

		data := []varieties.TradingVarieties{}

		q := db.Table(new(varieties.TradingVarieties))

		if search.Symbol != "" {
			q = q.Where("symbol like ?", "%"+search.Symbol+"%")
		}
		if search.Name != "" {
			q = q.Where("name like ?", "%"+search.Name+"%")
		}
		if search.Status != "" {
			q = q.Where("status = ?", search.Status)
		}

		cond := q.Conds()
		err := q.OrderBy("id desc").Limit(limit, offset).Find(&data)
		if err != nil {
			render(ctx, 1, err.Error(), 0, "")
			return
		}

		total, _ := db.Table(new(varieties.TradingVarieties)).And(cond).Count()
		if ctx.Query("api") == "1" {
			render(ctx, 0, "", int(total), data)
		} else {
			ctx.HTML(200, "tradingvarieties_list", gin.H{
				"searchParams": searchParams,
			})
		}
		return
	}
}
