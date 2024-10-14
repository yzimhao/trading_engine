package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderController struct{}

func NewOrderController() *OrderController {
	return &OrderController{}
}

// TODO implement
func (ctrl *OrderController) Create(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

// TODO implement
func (ctrl *OrderController) List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": "test",
	})
}

// TODO implement
func (ctrl *OrderController) Update(c *gin.Context) {}

// TODO implement
func (ctrl *OrderController) Delete(c *gin.Context) {}
