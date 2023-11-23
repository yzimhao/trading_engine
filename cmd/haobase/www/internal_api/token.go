package internal_api

import (
	"github.com/gin-gonic/gin"
	"github.com/gookit/goutil/arrutil"
	"github.com/yzimhao/trading_engine/types/token"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app/config"
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

	err := token.Set(req.Token, req.UserId, req.Ttl)
	if err != nil {
		utils.ResponseFailJson(ctx, err.Error())
		return
	}
	utils.ResponseOkJson(ctx, "")
}

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if len(config.App.Haobase.InternalApiAllowIp) > 0 && !arrutil.Contains(config.App.Haobase.InternalApiAllowIp, ip) {
			utils.ResponseFailJson(c, "非法IP")
			c.Abort()
			return
		}

		c.Next()
	}
}
