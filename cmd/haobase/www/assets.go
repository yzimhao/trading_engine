package www

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func assets_balance(ctx *gin.Context) {
	user_id := ctx.MustGet("user_id").(string)
	ss := ctx.Query("symbols")

	symbols := make([]string, 0)
	if ss != "" {
		symbols = strings.Split(ss, ",")
	}
	rows := assets.UserAssets(user_id, symbols)
	//todo 格式化资产数字

	for i, v := range rows {
		cfg, err := base.NewSymbols().Get(v.Symbol)
		if err != nil {
			app.Logger.Errorf("获取资产%s失败 %s", v.Symbol, err.Error())
			continue
		}

		rows[i].Total = utils.FormatDecimal(v.Total, cfg.ShowPrecision)
		rows[i].Freeze = utils.FormatDecimal(v.Freeze, cfg.ShowPrecision)
		rows[i].Available = utils.FormatDecimal(v.Available, cfg.ShowPrecision)
	}

	utils.ResponseOkJson(ctx, rows)
}
