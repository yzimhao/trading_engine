package orders

import (
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haoquote/quote/period"
	"github.com/yzimhao/trading_engine/utils"
)

type TradeStats struct {
	Symbol           string `json:"symbol"`
	TodayTradeQty    string `json:"today_trade_qty"`
	TodayTradeAmount string `json:"today_trade_amount"`
}

func NewTradeStats() []TradeStats {

	data := make([]TradeStats, 0)
	//从redis中的period_usdjpy_d1_1701532800_1701619199获取
	for _, v := range base.NewTSymbols().All() {
		stats, _ := period.GetTodyStats(v.Symbol)
		t := TradeStats{
			Symbol:           v.Symbol,
			TodayTradeQty:    utils.D(stats.Volume).String(),
			TodayTradeAmount: utils.D(stats.Amount).String(),
		}
		data = append(data, t)
	}

	return data
}
