package www

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func recharge_for_demo(ctx *gin.Context) {
	user_id := ctx.MustGet("user_id").(string)
	//自动为demo用户充值
	default_amount := "10000.00"
	all := base.NewSymbols().All()
	for _, item := range all {
		if assets.BalanceOfTotal(user_id, item.Symbol).Equal(decimal.Zero) {
			assets.SysDeposit(user_id, item.Symbol, default_amount, "sys.give:"+user_id)
		}
	}
	utils.ResponseOkJson(ctx, "")
}

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
