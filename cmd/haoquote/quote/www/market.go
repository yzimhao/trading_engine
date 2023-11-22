package www

import (
	"github.com/yzimhao/trading_engine/cmd/haobase/message"
	"github.com/yzimhao/trading_engine/cmd/haobase/message/ws"
	"github.com/yzimhao/trading_engine/cmd/haoquote/quote/period"
	"github.com/yzimhao/trading_engine/types"
	"github.com/yzimhao/trading_engine/utils"
)

// 24hr 价格变动情况
func market_24h(symbol string, last_price string) {
	var price string
	price, has := period.GetYesterdayClose(symbol)
	if !has {
		price, _ = period.GetTodayOpen(symbol)
	}

	to := types.MsgMarket24H.Format(map[string]string{
		"symbol": symbol,
	})

	message.Publish(ws.MsgBody{
		To: to,
		Response: ws.Response{
			Type: to,
			Body: map[string]any{
				"price_change": func() string {
					return utils.D(last_price).Sub(utils.D(price)).String()
				}(),
				"price_change_percent": func() string {
					c := utils.D(last_price).Sub(utils.D(price))
					if price != "" {
						return c.Div(utils.D(price)).Mul(utils.D("100")).StringFixed(2)
					}
					return "0"
				}(),
			},
		},
	})

}
