package token

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func Set(token string, user_id string, ttl int) {
	rdc := app.RedisPool().Get()
	defer rdc.Close()

	topic := tokenRedisTopic(token)
	rdc.Do("set", topic, user_id)
	rdc.Do("expire", topic, ttl)
}

func Get(original_token string) string {
	rdc := app.RedisPool().Get()
	defer rdc.Close()

	topic := tokenRedisTopic(original_token)
	user_id, err := redis.String(rdc.Do("get", topic))
	if err != nil {
		app.Logger.Errorf("从redis获取token信息出错 %s", err.Error())
	}
	return user_id
}

func tokenRedisTopic(token string) string {
	return fmt.Sprintf("user.token.%s", utils.Hash256(token))
}
