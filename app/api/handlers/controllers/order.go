package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	c.JSON(http.StatusOK, "ok")
}

// TODO implement
func (ctrl *OrderController) HistoryList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": "test",
	})
}

// TODO implement
func (ctrl *OrderController) UnfinishedList(c *gin.Context) {}

// TODO implement
func (ctrl *OrderController) TradeHistoryList(c *gin.Context) {}
