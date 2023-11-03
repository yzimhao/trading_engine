package orders

import "github.com/yzimhao/trading_engine/utils/app"

type TradeStats struct {
	TodayTradeQty    string `json:"today_trade_qty"`
	TodayTradeAmount string `json:"today_trade_amount"`
}

func TradeResultStats() TradeStats {
	db := app.Database().NewSession()
	defer db.Close()

	//todo
	return TradeStats{
		TodayTradeQty:    "1001",
		TodayTradeAmount: "10000.00",
	}
}
