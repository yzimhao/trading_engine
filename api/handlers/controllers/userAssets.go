package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserAssetsController struct{}

func NewUserAssetsController() *UserAssetsController {
	return &UserAssetsController{}
}

// TODO implement
func (ctrl *UserAssetsController) Create(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

// TODO implement
func (ctrl *UserAssetsController) List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": "test",
	})
}

// TODO implement
func (ctrl *UserAssetsController) Update(c *gin.Context) {}

// TODO implement
func (ctrl *UserAssetsController) Delete(c *gin.Context) {}
