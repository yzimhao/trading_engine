package orders

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/yzimhao/trading_engine/utils/app"
)

type LockType string

const (
	SettleLock LockType = "lock" //订单结算锁
)

func Lock(lt LockType, order_id string) {
	rdc := app.RedisPool().Get()
	defer rdc.Close()

	key := fmt.Sprintf("%s.%s", lt, order_id)
	rdc.Do("INCR", key)
}

func UnLock(lt LockType, order_id string) {
	rdc := app.RedisPool().Get()
	defer rdc.Close()

	key := fmt.Sprintf("%s.%s", lt, order_id)
	if _, err := rdc.Do("DECR", key); err != nil {
		app.Logger.Warnf("unlock %s err: %s", order_id, err.Error())
	}

	if n, _ := redis.Int64(rdc.Do("GET", key)); n == 0 {
		if _, err := rdc.Do("del", key); err != nil {
			app.Logger.Warnf("unlock %s fail err: %s", order_id, err.Error())
		}
	}
}

func GetLock(lt LockType, order_id string) int64 {
	rdc := app.RedisPool().Get()
	defer rdc.Close()

	key := fmt.Sprintf("%s.%s", lt, order_id)
	n, _ := redis.Int64(rdc.Do("GET", key))
	return n
}
