package view

import (
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haoadm/models"
	"github.com/yzimhao/trading_engine/cmd/haoadm/view/admin"
	"github.com/yzimhao/trading_engine/config"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func runModeCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if config.App.Main.Mode == app.ModeDemo.String() && config.App.Haoadm.Readonly {
			if ctx.Request.Method == "POST" {
				ctx.Abort()
				utils.ResponseFailJson(ctx, "demo模式，禁止修改数据")
				return
			}
		}
		ctx.Next()
	}
}

func recordLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		user_id := admin.GetLoginUserId(ctx)

		body, _ := ioutil.ReadAll(ctx.Request.Body)
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		models.NewAdminActivityLog(user_id, ctx.Request.Method, ctx.Request.RequestURI, string(body), ctx.ClientIP())
		ctx.Next()
	}
}
