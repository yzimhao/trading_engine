package www

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haoquote/quote/tradelog"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func trans_record(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	limit := utils.S2Int(ctx.Query("limit"))

	tsymbols := base.NewTradeSymbol()
	info, err := tsymbols.Get(symbol)
	if err != nil {
		utils.ResponseFailJson(ctx, err.Error())
		return
	}

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

	price, qty := info.PricePrecision, info.QtyPrecision

	for i, v := range rows {
		rows[i].TradePrice = utils.NumberFix(v.TradePrice, int(price))
		rows[i].TradeAmount = utils.NumberFix(v.TradeAmount, int(price))
		rows[i].TradeQuantity = utils.NumberFix(v.TradeQuantity, int(qty))
	}

	utils.ResponseOkJson(ctx, rows)
}
