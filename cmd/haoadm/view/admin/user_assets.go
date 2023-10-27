package admin

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func AssetsList(ctx *gin.Context) {
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

		data := []assets.Assets{}

		q := db.Table(new(assets.Assets))

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

		total, _ := q.And(cond).Count()
		if ctx.Query("api") == "1" {
			render(ctx, 0, "", int(total), data)
		} else {
			ctx.HTML(200, "user_assets", gin.H{
				"searchParams": searchParams,
			})
		}
		return
	}
}
