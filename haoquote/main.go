package haoquote

import (
	"context"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/yzimhao/trading_engine/haoquote/tradelog"
	"github.com/yzimhao/trading_engine/haoquote/www"
	"xorm.io/xorm"
)

func Start(ctx *context.Context, rc *redis.Client, db *xorm.Engine) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	init_symbols_quote()
	//http
	tradelog.Init(rc, db)
	www.Run(rc, db)
	wg.Wait()
}

func init_symbols_quote() {
	symbols := viper.GetStringMap("symbol")

	for k, attr := range symbols {
		symbol := strings.ToLower(k)
		price_digit := attr.(map[string]any)["price_digit"].(int64)
		qty_digit := attr.(map[string]any)["qty_digit"].(int64)
		go tradelog.Monitor(symbol, price_digit, qty_digit)
	}
}
