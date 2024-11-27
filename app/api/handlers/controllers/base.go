package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/app/api/handlers/common"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type BaseController struct {
	logger *zap.Logger
}

type inBaseContext struct {
	fx.In
	Logger *zap.Logger
}

func NewBaseController(in inBaseContext) *BaseController {
	return &BaseController{logger: in.Logger}
}

// @Summary ping
// @Description test if the server is running
// @ID v1.ping
// @Tags base
// @Accept json
// @Produce json
// @Success 200 {string} any
// @Router /api/v1/ping [get]
func (ctrl *BaseController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (ctrl *BaseController) Time(c *gin.Context) {
	common.ResponseOK(c, gin.H{
		"time": time.Now().Nanosecond(),
	})
}

// @Summary exchange info
// @Description get exchange info
// @ID v1.base.exchange_info
// @Tags base
// @Accept json
// @Produce json
// @Query param string true "symbol"
// @Success 200 {string} any
// @Router /api/v1/base/exchange_info [get]
func (ctrl *BaseController) ExchangeInfo(c *gin.Context) {
	//TODO
	common.ResponseOK(c, gin.H{
		"data": "not implemented",
	})
}
