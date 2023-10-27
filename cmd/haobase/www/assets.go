package www

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

// 用户资产余额返回
type response_assets struct {
	Symbol     string     `json:"symbol"`
	Total      string     `json:"total"`
	Freeze     string     `json:"freeze"`
	Available  string     `json:"avail"`
	UpdateTime utils.Time `json:"update_time"`
}

func assets_balance(ctx *gin.Context) {
	user_id := ctx.MustGet("user_id").(string)
	ss := ctx.Query("symbols")

	symbols := make([]string, 0)
	if ss != "" {
		symbols = strings.Split(ss, ",")
	}
	rows := assets.UserAssets(user_id, symbols)

	data := make([]response_assets, 0)

	for _, v := range rows {
		cfg, err := base.NewSymbols().Get(v.Symbol)
		if err != nil {
			app.Logger.Errorf("获取资产%s失败 %s", v.Symbol, err.Error())
			continue
		}

		item := response_assets{
			Symbol:    v.Symbol,
			Total:     utils.FormatDecimal(v.Total, cfg.ShowPrecision),
			Freeze:    utils.FormatDecimal(v.Freeze, cfg.ShowPrecision),
			Available: utils.FormatDecimal(v.Available, cfg.ShowPrecision),
		}

		data = append(data, item)
	}

	utils.ResponseOkJson(ctx, data)
}
