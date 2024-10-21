package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/app/api/handlers/common"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MarketController struct {
	logger *zap.Logger
}

type inMarketContext struct {
	fx.In
	Logger *zap.Logger
}

func NewMarketController(in inMarketContext) *MarketController {
	return &MarketController{
		logger: in.Logger,
	}
}

// @Summary depth
// @Description get depth
// @ID v1.market.depth
// @Tags market
// @Accept json
// @Produce json
// @Param symbol query string true "symbol"
// @Param limit query int false "limit"
// @Success 200 {string} any
// @Router /api/v1/market/depth [get]
func (ctrl *MarketController) Depth(c *gin.Context) {
	common.ResponseOK(c, gin.H{})
}

// @Summary trades
// @Description 获取近期成交记录
// @ID v1.market.trades
// @Tags market
// @Accept json
// @Produce json
// @Param symbol query string true "symbol"
// @Param limit query int false "limit"
// @Success 200 {string} any
// @Router /api/v1/market/trades [get]
func (ctrl *MarketController) Trades(c *gin.Context) {
	common.ResponseOK(c, gin.H{})
}

// @Summary klines
// @Description 获取K线数据
// @ID v1.market.klines
// @Tags market
// @Accept json
// @Produce json
// @Param symbol query string true "symbol"
// @Param limit query int false "limit"
// @Success 200 {string} any
// @Router /api/v1/market/klines [get]
func (ctrl *MarketController) Klines(c *gin.Context) {
	common.ResponseOK(c, gin.H{})
}
