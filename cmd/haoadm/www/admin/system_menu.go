package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haoadm/models"
)

func SystemMenu(ctx *gin.Context) {
	s := models.SystemMenu{}
	data := s.GetSystemInit()
	ctx.JSON(200, data)
}
