package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/app/api/handlers/common"
	"go.uber.org/fx"
)

type OrderController struct{}

type inOrderContext struct {
	fx.In
}

func NewOrderController(in inOrderContext) *OrderController {
	return &OrderController{}
}

// @Summary create order
// @Description create order
// @ID v1.order.create
// @Tags order
// @Accept json
// @Produce json
// @Success 200 {string} any
// @Router /api/v1/order/create [post]
func (ctrl *OrderController) Create(c *gin.Context) {
	common.ResponseOK(c, "ok")
}

// @Summary history list
// @Description history list
// @ID v1.order.history
// @Tags order
// @Accept json
// @Produce json
// @Success 200 {string} any
// @Router /api/v1/order/history [get]
func (ctrl *OrderController) HistoryList(c *gin.Context) {
	common.ResponseOK(c, "test")
}

// @Summary unfinished list
// @Description unfinished list
// @ID v1.order.unfinished
// @Tags order
// @Accept json
// @Produce json
// @Success 200 {string} any
// @Router /api/v1/order/unfinished [get]
func (ctrl *OrderController) UnfinishedList(c *gin.Context) {
	common.ResponseOK(c, "test")
}

// @Summary trade history list
// @Description trade history list
// @ID v1.order.trade_history
// @Tags order
// @Accept json
// @Produce json
// @Success 200 {string} any
// @Router /api/v1/order/trade/history [get]
func (ctrl *OrderController) TradeHistoryList(c *gin.Context) {
	common.ResponseOK(c, "test")
}
