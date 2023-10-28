package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haoadm/models"
)

func SystemMenu(ctx *gin.Context) {
	if ctx.Query("v") == "v2" {
		s := models.SystemMenu{}
		data := s.GetV1SystemInit()
		ctx.JSON(200, data)
	} else {
		s := models.SystemMenu{}
		data := s.GetV1SystemInit()
		ctx.JSON(200, data)
	}

}
