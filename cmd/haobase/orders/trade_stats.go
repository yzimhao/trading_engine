package orders

import (
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haoquote/quote/period"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

type TradeStats struct {
	TodayTradeQty    string `json:"today_trade_qty"`
	TodayTradeAmount string `json:"today_trade_amount"`
}

func NewTradeStats() TradeStats {
	db := app.Database().NewSession()
	defer db.Close()

	ts := TradeStats{
		TodayTradeQty:    "0",
		TodayTradeAmount: "0.00",
	}
	ts.stats()
	return ts
}

func (ts *TradeStats) stats() {
	//从redis中的period_usdjpy_d1_1701532800_1701619199获取
	for _, v := range base.NewTSymbols().All() {
		data, _ := period.GetTodyStats(v.Symbol)
		ts.TodayTradeQty = utils.D(ts.TodayTradeQty).Add(utils.D(data.Volume)).String()
		ts.TodayTradeAmount = utils.D(ts.TodayTradeAmount).Add(utils.D(data.Amount)).String()
	}
}
