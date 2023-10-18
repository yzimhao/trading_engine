package haoquote

import (
	"strings"

	"github.com/spf13/viper"
	"github.com/yzimhao/trading_engine/haoquote/tradelog"
	"github.com/yzimhao/trading_engine/haoquote/www"
	"github.com/yzimhao/trading_engine/utils/filecache"
)

func Run() {
	init_symbols_quote()
	tradelog.Init()
	www.Run()
}

func init_symbols_quote() {
	symbols := viper.GetStringMap("symbol")
	filecache.NewStorage(viper.GetString("haoquote.storage_path"), 1)

	for k, attr := range symbols {
		symbol := strings.ToLower(k)
		price_digit := attr.(map[string]any)["price_digit"].(int64)
		qty_digit := attr.(map[string]any)["qty_digit"].(int64)
		go tradelog.Monitor(symbol, price_digit, qty_digit)
	}
}
