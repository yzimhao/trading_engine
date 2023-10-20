package internal_api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

type req_settoken_args struct {
	UserId string `json:"user_id" binding:"required"`
	Token  string `json:"token" binding:"required"`
	Ttl    int    `json:"ttl" binding:"required"`
}

func SetToken(ctx *gin.Context) {
	var req req_settoken_args
	if err := ctx.BindJSON(&req); err != nil {
		utils.ResponseFailJson(ctx, err.Error())
		return
	}
	UpdateRedisToken(req)
	utils.ResponseOkJson(ctx, "")
}

func UpdateRedisToken(req req_settoken_args) {
	rdc := app.RedisPool().Get()
	defer rdc.Close()

	topic := tokenRedisTopic(req.Token)
	rdc.Do("set", topic, req.UserId)
	rdc.Do("expire", topic, req.Ttl)
}

func GetUserIdFromToken(original_token string) string {
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
