package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/yzimhao/trading_engine/v2/internal/persistence"
)

type UserAssetsController struct {
	repo   persistence.AssetsRepository
	logger *zap.Logger
}

func NewUserAssetsController(repo persistence.AssetsRepository, logger *zap.Logger) *UserAssetsController {

	return &UserAssetsController{repo: repo, logger: logger}
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

func (ctrl *UserAssetsController) Query(c *gin.Context) {

	symbol := c.Param("symbol")

	ctx := context.Background()
	err := ctrl.repo.Despoit(ctx, "user1", symbol, "1")
	if err != nil {
		ctrl.logger.Error("Query", zap.Error(err))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"symbol": symbol,
			"error":  err,
		},
	})
}
