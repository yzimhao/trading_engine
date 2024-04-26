package admin

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

type assetsSearch struct {
	Symbol     string `json:"symbol"`
	UserId     string `json:"user_id"`
	Status     string `json:"status"`
	BusinessId string `json:"business_id"`
	ChangeType string `json:"change_type"`
}

func AssetsList(ctx *gin.Context) {
	db := app.Database().NewSession()
	defer db.Close()

	if ctx.Request.Method == "GET" {
		page := utils.S2Int(ctx.Query("page"))
		limit := utils.S2Int(ctx.Query("limit"))
		searchParams := ctx.Query("searchParams")

		var search assetsSearch
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
			q = q.Where("symbol = ?", search.Symbol)
		}

		if search.UserId != "" {
			q = q.Where("user_id = ?", search.UserId)
		}

		cond := q.Conds()
		err := q.OrderBy("id desc").Limit(limit, offset).Find(&data)
		if err != nil {
			render(ctx, 1, err.Error(), 0, "")
			return
		}

		total, _ := db.Table(new(assets.Assets)).And(cond).Count()
		if ctx.Query("api") == "1" {
			render(ctx, 0, "", int(total), data)
		} else {
			ctx.HTML(200, "user_assets", gin.H{
				"search":      search,
				"all_symbols": base.NewSymbols().All(),
			})
		}
		return
	}

}

func AssetsModify(ctx *gin.Context) {
	//assetsArgs := ctx.Query("assetsArgs")
	//
	//var asset assetsVo
	//json.Unmarshal([]byte(assetsArgs), &asset)
	//avail := asset.Avail
	//userId := asset.UserId
	//symbol := asset.Symbol
	ctx.HTML(200, "user_assets_modify", gin.H{
		//"avail":  avail,
		//"userId": userId,
		//"symbol": symbol,
	})
}

func AddUserAssets(ctx *gin.Context) {
	symbol := ctx.PostForm("symbol")
	amount := ctx.PostForm("amount")
	userId := ctx.PostForm("userId")
	_, err := assets.SysDeposit(userId, symbol, amount, "sys.give:"+userId)
	if err != nil {
		utils.ResponseFailJson(ctx, err.Error())
		return
	}
	utils.ResponseOkJson(ctx, gin.H{})
	return
}

func DecreaseUserAssets(ctx *gin.Context) {
	symbol := ctx.PostForm("symbol")
	amount := ctx.PostForm("amount")
	userId := ctx.PostForm("userId")
	_, err := assets.SysWithdraw(userId, symbol, amount, "sys.decrease:"+userId)
	if err != nil {
		utils.ResponseFailJson(ctx, err.Error())
		return
	}
	utils.ResponseOkJson(ctx, gin.H{})
	return
}
