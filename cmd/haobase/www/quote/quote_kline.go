package quote

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haoquote/quote/period"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func KLine(ctx *gin.Context) {
	interval := ctx.Query("interval")
	limit := utils.S2Int(ctx.Query("limit"))
	symbol := ctx.Query("symbol")

	tsymbols := base.NewTradeSymbol()
	info, err := tsymbols.Get(symbol)
	if err != nil {
		utils.ResponseFailJson(ctx, err.Error())
		return
	}

	if limit > 1000 {
		limit = 1000
	}
	if limit <= 0 {
		limit = 500
	}

	row := period.Period{
		Symbol:   symbol,
		Interval: period.PeriodType(interval),
	}

	// [
	//     [
	//       1499040000000,      // k线开盘时间
	//       "0.01634790",       // 开盘价
	//       "0.80000000",       // 最高价
	//       "0.01575800",       // 最低价
	//       "0.01577100",       // 收盘价(当前K线未结束的即为最新价)
	//       "148976.11427815",  // 成交量
	//       1499644799999,      // k线收盘时间
	//       "2434.19055334",    // 成交额
	//       308,                // 成交笔数
	//       "1756.87402397",    // 主动买入成交量
	//       "28.46694368",      // 主动买入成交额
	//       "17928899.62484339" // 请忽略该参数
	//     ]
	//   ]

	db := app.Database().NewSession()
	defer db.Close()

	var rows []period.Period
	db.Table(row.TableName()).OrderBy("open_at desc").Limit(limit).Find(&rows)

	pd, qd := info.PricePrecision, info.QtyPrecision

	data := make([][6]any, 0)
	for _, v := range rows {
		data = append(data, [6]any{
			time.Time(v.OpenAt).UnixMilli(),
			utils.NumberFix(v.Open, int(pd)),
			utils.NumberFix(v.High, int(pd)),
			utils.NumberFix(v.Low, int(pd)),
			utils.NumberFix(v.Close, int(pd)),
			utils.NumberFix(v.Volume, int(qd)),
		})
	}

	utils.ResponseOkJson(ctx, data)
}
