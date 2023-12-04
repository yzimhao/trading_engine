package period

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/yzimhao/trading_engine/utils/app"
)

type periodCachekey string

const (
	periodBucket string         = "period"
	periodKey    periodCachekey = "period_%s_%s_%d_%d" //"period_usdjpy_mn_1695571200_1696175999"
)

func (p periodCachekey) Format(pt PeriodType, symbol string, st, et int64) periodCachekey {
	v := fmt.Sprintf(string(p), symbol, pt, st, et)
	return periodCachekey(v)
}

func (p periodCachekey) set(value []byte, ttl int64) {
	rc := app.RedisPool().Get()
	defer rc.Close()

	rc.Do("set", p, value)
	rc.Do("expire", p, ttl)
}

func (p periodCachekey) get() ([]byte, error) {
	rc := app.RedisPool().Get()
	defer rc.Close()

	return redis.Bytes(rc.Do("get", p))
}

func GetYesterdayClose(symbol string) (string, bool) {
	now := time.Now()

	//获取昨天的收盘价，如果没有则获取今天的开盘价
	st, et := get_start_end_time(now.AddDate(0, 0, -1), PERIOD_D1)
	key := periodKey.Format(PERIOD_D1, symbol, st.Unix(), et.Unix())
	cache_data, err := key.get()
	if err != nil {
		return "", false
	}

	var data Period
	if err := json.Unmarshal(cache_data, &data); err != nil {
		return "", false
	}

	return data.Close, true
}

func GetTodyStats(symbol string) (Period, error) {
	now := time.Now()

	st, et := get_start_end_time(now, PERIOD_D1)
	key := periodKey.Format(PERIOD_D1, symbol, st.Unix(), et.Unix())
	cache_data, err := key.get()

	var data Period

	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(cache_data, &data); err != nil {
		return data, err
	}
	return data, nil
}
