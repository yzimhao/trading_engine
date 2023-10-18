package haotrader

import (
	"sync"
	"time"

	"github.com/gookit/goutil/arrutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/utils/app"
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
	wg.Add(1)
	logrus.Info("启动撮合程序成功! 如需帮助请参考: https://github.com/yzimhao/trading_engine")
	init_symbols_tengine()
	wg.Wait()
}

func init_symbols_tengine() {
	local_config_symbols := app.CstringSlice("local.symbols")
	db_symbols := base.NewTSymbols().All()
	for _, item := range db_symbols {
		if len(local_config_symbols) > 0 && arrutil.Contains(local_config_symbols, item.Symbol) || len(local_config_symbols) == 0 {
			teps[item.Symbol] = NewTengine(item.Symbol, item.PricePrecision, item.QtyPrecision)
		}
	}
}
