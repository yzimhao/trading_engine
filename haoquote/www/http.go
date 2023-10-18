package www

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/yzimhao/trading_engine/cmd/haobase/message"
	"github.com/yzimhao/trading_engine/cmd/haobase/message/ws"
	_ "github.com/yzimhao/trading_engine/docs/api" // main 文件中导入 docs 包
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"

	"github.com/yzimhao/trading_engine/haoquote/period"
	"github.com/yzimhao/trading_engine/haoquote/tradelog"
)

var ()

func Run() {
	sub_symbol_depth()
	http_start(viper.GetString("haoquote.http.host"))
}

func http_start(addr string) {
	if viper.GetBool("haoquote.http.debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	logrus.Infof("HTTP服务监听: %s", addr)
	web_router(router)
	router.Run(addr)
}

func web_router(router *gin.Engine) {
	//websokect服务放在这个quote里
	ws.NewHub()
	message.Subscribe()

	router.GET("/quote/ws", func(ctx *gin.Context) {
		ws.M.ServeWs(ctx)
	})

	api := router.Group("/api/v1/quote")
	{
		api.Use(utils.CorsMiddleware())
		api.GET("/depth", symbol_depth)
		api.GET("/trans/record", trans_record)
		api.GET("/kline", kline)
		api.GET("/system", system_info)
	}
}

func symbol_depth(ctx *gin.Context) {
	limit := utils.S2Int(ctx.Query("limit"))
	symbol := strings.ToLower(ctx.Query("symbol"))

	if _, ok := symbols_depth.data[symbol]; !ok {
		utils.ResponseFailJson(ctx, "invalid symbol")
		return
	}

	if limit <= 0 || limit > 100 {
		limit = 10
	}

	utils.ResponseOkJson(ctx, gin.H{
		"asks": symbols_depth.limit("asks", symbol, limit),
		"bids": symbols_depth.limit("bids", symbol, limit),
	})
}

func trans_record(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	limit := utils.S2Int(ctx.Query("limit"))

	rows := make([]tradelog.TradeLog, 0)
	tl := tradelog.TradeLog{
		Symbol: symbol,
	}

	db := app.Database().NewSession()
	defer db.Close()

	db.Table(tl.TableName()).OrderBy("trade_at desc, id desc").Limit(limit).Find(&rows)

	// [
	//     {
	//         "id": 28457,
	//         "price": "4.00000100",
	//         "qty": "12.00000000",
	//         "time": 1499865549590, // 交易成交时间, 和websocket中的T一致.
	//         "isBuyerMaker": true,
	//         "isBestMatch": true
	//     }
	// ]

	price, qty := symbols_depth.get_digit(symbol)

	for i, v := range rows {
		rows[i].TradePrice = utils.NumberFix(v.TradePrice, int(price))
		rows[i].TradeAmount = utils.NumberFix(v.TradeAmount, int(price))
		rows[i].TradeQuantity = utils.NumberFix(v.TradeQuantity, int(qty))
	}

	utils.ResponseOkJson(ctx, rows)
}

func kline(ctx *gin.Context) {
	interval := ctx.Query("interval")
	limit := utils.S2Int(ctx.Query("limit"))
	symbol := ctx.Query("symbol")

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

	pd, qd := symbols_depth.get_digit(symbol)

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

func system_info(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"version": app.Version,
		"build":   app.Build,
	})
}
