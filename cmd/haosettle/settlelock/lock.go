package settlelock

import (
	"sync"

	"github.com/gomodule/redigo/redis"
	"github.com/yzimhao/trading_engine/types/redisdb"
	"github.com/yzimhao/trading_engine/utils/app"
)

var (
	m sync.Mutex
)

func Lock(order_id ...string) {
	m.Lock()
	defer m.Unlock()

	rdc := app.RedisPool().Get()
	defer rdc.Close()

	for _, oid := range order_id {
		key := redisdb.OrderLock.Format(redisdb.Replace{"order_id": oid})
		if _, err := rdc.Do("INCR", key); err != nil {
			app.Logger.Errorf("lock order %s err: %s", oid, err.Error())
		}
	}

}

func UnLock(order_id ...string) {
	m.Lock()
	defer m.Unlock()

	rdc := app.RedisPool().Get()
	defer rdc.Close()

	for _, oid := range order_id {
		key := redisdb.OrderLock.Format(redisdb.Replace{"order_id": oid})

		if _, err := rdc.Do("DECR", key); err != nil {
			app.Logger.Errorf("unlock %s err: %s", oid, err.Error())
		}

		if n, _ := redis.Int64(rdc.Do("GET", key)); n == 0 {
			if _, err := rdc.Do("del", key); err != nil {
				app.Logger.Errorf("unlock %s fail err: %s", oid, err.Error())
			}
		}
	}

}

func GetLock(order_id string) int64 {
	rdc := app.RedisPool().Get()
	defer rdc.Close()

	key := redisdb.OrderLock.Format(redisdb.Replace{"order_id": order_id})

	n, _ := redis.Int64(rdc.Do("GET", key))
	return n
}
