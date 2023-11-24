package token

import (
	"fmt"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func Set(token string, user_id string, ttl int) error {
	rdc := app.RedisPool().Get()
	defer rdc.Close()

	if strings.Contains(token, ".") {
		return fmt.Errorf("token contains `.`")
	}

	topic := tokenRedisTopic(token)
	if _, err := rdc.Do("set", topic, user_id); err != nil {
		return err
	}
	if _, err := rdc.Do("expire", topic, ttl); err != nil {
		return err
	}

	return nil
}

func Get(original_token string) string {
	rdc := app.RedisPool().Get()
	defer rdc.Close()

	topic := tokenRedisTopic(original_token)
	user_id, err := redis.String(rdc.Do("get", topic))
	if err != nil {
		app.Logger.Infof("从redis获取token信息出错 %s", err.Error())
	}
	return user_id
}

func tokenRedisTopic(token string) string {
	return fmt.Sprintf("user.token.%s", utils.Hash256(token))
}
