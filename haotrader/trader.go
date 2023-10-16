package haotrader

import (
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/utils/filecache"
)

var (
	wg sync.WaitGroup

	localdb *filecache.Storage
	teps    map[string]*trading_core.TradePair
)

func Run() {
	teps = make(map[string]*trading_core.TradePair)
	localdb = filecache.NewStorage(viper.GetString("haotrader.storage_path"), time.Duration(10))
	defer localdb.Close()

	wg = sync.WaitGroup{}
	//todo wg.done()

	logrus.Info("启动撮合程序成功! 如需帮助请参考: https://github.com/yzimhao/trading_engine")
	init_symbols_tengine()

	if viper.GetString("haotrader.http.host") != "" {
		go http_start(viper.GetString("haotrader.http.host"))
	}
	wg.Wait()
}

func init_symbols_tengine() {
	symbols := viper.GetStringMap("symbol")

	for k, attr := range symbols {
		wg.Add(1)
		symbol := strings.ToLower(k)
		price_digit := attr.(map[string]any)["price_digit"].(int64)
		qty_digit := attr.(map[string]any)["qty_digit"].(int64)
		teps[symbol] = NewTengine(symbol, int(price_digit), int(qty_digit))
	}
}
