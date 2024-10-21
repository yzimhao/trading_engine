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

// TODO implement
func (ctrl *OrderController) Create(c *gin.Context) {
	common.ResponseOK(c, "ok")
}

// TODO implement
func (ctrl *OrderController) HistoryList(c *gin.Context) {
	common.ResponseOK(c, "test")
}

// TODO implement
func (ctrl *OrderController) UnfinishedList(c *gin.Context) {}

// TODO implement
func (ctrl *OrderController) TradeHistoryList(c *gin.Context) {}
