package orders

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/types/redisdb"
	"github.com/yzimhao/trading_engine/utils/app"
)

func generate_order_id_by_side(side trading_core.OrderSide) string {
	if side == trading_core.OrderSideSell {
		return generate_order_id("A")
	} else {
		return generate_order_id("B")
	}
}

// todo 调整order_id的长度
func generate_order_id(prefix string) string {
	prefix = strings.ToUpper(prefix)
	t := time.Now().Format("060102150405")
	rn := rand.Intn(9999)
	return fmt.Sprintf("%s%s%04d", prefix, t, rn)
}

func push_new_order_to_redis(symbol string, data []byte) {
	topic := redisdb.NewOrderQueue.Format(redisdb.Replace{"symbol": symbol})
	app.Logger.Infof("推送新订单%s: %s", topic, data)

	rdc := app.RedisPool().Get()
	defer rdc.Close()
	if _, err := rdc.Do("RPUSH", topic, data); err != nil {
		app.Logger.Errorf("RPUSH %s err: %s", topic, err.Error())
	}
}
