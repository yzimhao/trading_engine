package admin

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func AssetsFreezeList(ctx *gin.Context) {
	db := app.Database().NewSession()
	defer db.Close()

	if ctx.Request.Method == "GET" {
		page := utils.S2Int(ctx.Query("page"))
		limit := utils.S2Int(ctx.Query("limit"))
		searchParams := ctx.Query("searchParams")

		var search assetsSearch
		json.Unmarshal([]byte(searchParams), &search)
		app.Logger.Debugf("search: %#v args: %s", search, searchParams)

		if page <= 0 {
			page = 1
		}
		if limit <= 0 {
			limit = 10
		}
		offset := (page - 1) * limit

		data := []assets.AssetsFreeze{}
		tb := assets.AssetsFreeze{Symbol: search.Symbol}
		q := db.Table(tb.TableName())

		if search.Status != "" {
			q = q.Where("status = ?", search.Status)
		}
		if search.UserId != "" {
			q = q.Where("user_id = ?", search.UserId)
		}
		if search.BusinessId != "" {
			q = q.Where("business_id=?", search.BusinessId)
		}

		cond := q.Conds()
		err := q.OrderBy("id desc").Limit(limit, offset).Find(&data)
		if err != nil {
			render(ctx, 1, err.Error(), 0, "")
			return
		}

		total, _ := db.Table(tb.TableName()).And(cond).Count()
		if ctx.Query("api") == "1" {
			render(ctx, 0, "", int(total), data)
		} else {
			ctx.HTML(200, "user_assets_freeze", gin.H{
				"search": search,
			})
		}
		return
	}
}
