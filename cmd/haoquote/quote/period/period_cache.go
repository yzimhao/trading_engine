package period

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/yzimhao/trading_engine/types/redisdb"
	"github.com/yzimhao/trading_engine/utils/app"
)

func formatKLineKey(pt PeriodType, symbol string, st, et int64) string {
	k := redisdb.KLinePeriod.Format(redisdb.Replace{"symbol": symbol, "period": string(pt), "start_time": fmt.Sprintf("%d", st), "end_time": fmt.Sprintf("%d", et)})
	return k
}

func setKLinePeriod(key string, value []byte, ttl int64) {
	rc := app.RedisPool().Get()
	defer rc.Close()

	if _, err := rc.Do("set", key, value); err != nil {
		app.Logger.Errorf("set kline period key %s error: %s", key, err)
	}
	if _, err := rc.Do("expire", key, ttl); err != nil {
		app.Logger.Errorf("set kline period key %s expire error: %s", key, err)
	}
}

func getKLinePeriod(key string) ([]byte, error) {
	rc := app.RedisPool().Get()
	defer rc.Close()

	return redis.Bytes(rc.Do("get", key))
}

func GetYesterdayClose(symbol string) (string, bool) {
	now := time.Now()

	//获取昨天的收盘价，如果没有则获取今天的开盘价
	st, et := parse_start_end_time(now.AddDate(0, 0, -1), PERIOD_D1)
	key := formatKLineKey(PERIOD_D1, symbol, st.Unix(), et.Unix())
	cache_data, err := getKLinePeriod(key)
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

	st, et := parse_start_end_time(now, PERIOD_D1)
	key := formatKLineKey(PERIOD_D1, symbol, st.Unix(), et.Unix())
	cache_data, err := getKLinePeriod(key)

	var data Period

	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(cache_data, &data); err != nil {
		return data, err
	}
	return data, nil
}
