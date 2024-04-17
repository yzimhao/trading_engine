package quote

import (
	"time"

	"github.com/gookit/goutil/arrutil"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haoquote/quote/tradelog"
	"github.com/yzimhao/trading_engine/cmd/haoquote/quote/www"
	"github.com/yzimhao/trading_engine/config"
	"github.com/yzimhao/trading_engine/utils/app/keepalive"
)

func Run() {
	go func() {
		for {
			init_symbols()
			time.Sleep(time.Second * 5)
		}
	}()
	www.Run()
}

func init_symbols() {
	local_config_symbols := config.App.Local.Symbols
	db_symbols := base.NewTradeSymbol().All()
	for _, item := range db_symbols {
		if len(local_config_symbols) > 0 && arrutil.Contains(local_config_symbols, item.Symbol) || len(local_config_symbols) == 0 {
			if !keepalive.HasExtrasKeyValue("quote.symbols", item.Symbol) {
				go tradelog.Monitor(item.Symbol, int64(item.PricePrecision), int64(item.QtyPrecision))
				keepalive.SetExtras("quote.symbols", item.Symbol)
			}

		}
	}
}
