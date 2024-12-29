package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/app/common"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type BaseController struct {
	logger       *zap.Logger
	tradeVariety persistence.TradeVarietyRepository
}

type inBaseContext struct {
	fx.In
	Logger       *zap.Logger
	TradeVariety persistence.TradeVarietyRepository
}

func NewBaseController(in inBaseContext) *BaseController {
	return &BaseController{
		logger:       in.Logger,
		tradeVariety: in.TradeVariety,
	}
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
	symbol := strings.ToUpper(c.Query("symbol"))

	tradeVariety, err := ctrl.tradeVariety.FindBySymbol(c, symbol)
	if err != nil {
		common.ResponseError(c, err)
		return
	}

	common.ResponseOK(c, tradeVariety)
}
