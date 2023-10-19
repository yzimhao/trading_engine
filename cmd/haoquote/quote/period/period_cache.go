package period

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/yzimhao/trading_engine/utils/filecache"
)

type periodCachekey string

const (
	periodBucket string         = "period"
	periodKey    periodCachekey = "period_%s_%s_%d_%d" //"period_usdjpy_mn_1695571200_1696175999"
)

func (c periodCachekey) Format(pt PeriodType, symbol string, st, et int64) string {
	return fmt.Sprintf(string(c), symbol, pt, st, et)
}

func newCache() *filecache.Storage {
	return filecache.NewStorage(viper.GetString("haoquote.cache"), 1)
}

func GetYesterdayClose(symbol string) (string, bool) {
	now := time.Now()
	cache := newCache()

	//获取昨天的收盘价，如果没有则获取今天的开盘价
	st, et := get_start_end_time(now.AddDate(0, 0, -1), PERIOD_D1)
	key := periodKey.Format(PERIOD_D1, symbol, st.Unix(), et.Unix())
	cache_data, has := cache.Get(periodBucket, key)

	var data Period
	json.Unmarshal(cache_data, &data)
	return data.Close, has
}

func GetTodayOpen(symbol string) (string, bool) {
	now := time.Now()
	cache := newCache()

	st, et := get_start_end_time(now, PERIOD_D1)
	key := periodKey.Format(PERIOD_D1, symbol, st.Unix(), et.Unix())
	cache_data, has := cache.Get(periodBucket, key)

	var data Period
	json.Unmarshal(cache_data, &data)
	return data.Open, has
}
