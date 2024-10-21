package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
)

type UserAssetsController struct {
	repo   persistence.AssetRepository
	logger *zap.Logger
}

func NewUserAssetsController(repo persistence.AssetRepository, logger *zap.Logger) *UserAssetsController {
	return &UserAssetsController{repo: repo, logger: logger}
}

type DespoitRequest struct {
	UserId string `json:"user_id"`
	Symbol string `json:"symbol"`
	Amount string `json:"amount"`
}

// @Summary asset despoit
// @Description despoit an asset
// @ID v1.asset.despoit
// @Tags asset
// @Accept json
// @Produce json
// @Param body body DespoitRequest true "despoit request"
// @Success 200 {string} order_id
// @Router /api/v1/asset/despoit [post]
func (ctrl *UserAssetsController) Despoit(c *gin.Context) {
	var req DespoitRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	order_id, err := ctrl.repo.Despoit(ctx, req.UserId, req.Symbol, req.Amount)
	if err != nil {
		ctrl.logger.Error("Despoit", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":   true,
		"data": order_id,
	})

}

type WithdrawRequest struct {
	UserId string `json:"user_id"`
	Symbol string `json:"symbol"`
	Amount string `json:"amount"`
}

// @Summary asset withdraw
// @Description withdraw an asset
// @ID v1.asset.withdraw
// @Tags asset
// @Accept json
// @Produce json
// @Param body body WithdrawRequest true "withdraw request"
// @Success 200 {string} order_id
// @Router /api/v1/asset/withdraw [post]
func (ctrl *UserAssetsController) Withdraw(c *gin.Context) {
	var req WithdrawRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	order_id, err := ctrl.repo.Withdraw(ctx, req.UserId, req.Symbol, req.Amount)
	if err != nil {
		ctrl.logger.Error("Withdraw", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":   true,
		"data": order_id,
	})
}

// @Summary get wallet asset
// @Description get an asset balance
// @ID v1.asset.query
// @Tags asset
// @Accept json
// @Produce json
// @Param symbol path string true "symbol"
// @Query userId query string true "userId测试用参数"
// @Success 200 {object} entities.Asset
// @Router /api/v1/asset/{symbol} [get]
func (ctrl *UserAssetsController) Query(c *gin.Context) {

	symbol := c.Param("symbol")
	userId := c.Query("userId")

	var asset *entities.Asset

	ctx := context.Background()
	asset, err := ctrl.repo.FindOne(ctx, userId, symbol)
	if err != nil {
		ctrl.logger.Error("Query", zap.Error(err))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": asset,
		"ok":   true,
	})
}

type TransferRequest struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Symbol string `json:"symbol"`
	Amount string `json:"amount"`
}

// @Summary asset transfer
// @Description transfer an asset
// @ID v1.asset.transfer
// @Tags asset
// @Accept json
// @Produce json
// @Param symbol path string true "symbol"
// @Param body body TransferRequest true "transfer request"
// @Router /api/v1/asset/transfer/{symbol} [post]
func (ctrl *UserAssetsController) Transfer(c *gin.Context) {

	var req TransferRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	err = ctrl.repo.Transfer(ctx, req.From, req.To, req.Symbol, req.Amount)
	if err != nil {
		ctrl.logger.Error("Query", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}

// @Summary get asset history
// @Description get an asset history
// @ID v1.asset.history
// @Tags asset
// @Accept json
// @Produce json
// @Param symbol path string true "symbol"
// @Success 200 {object} []entities.Asset
// @Router /api/v1/asset/{symbol}/history [get]
func (ctrl *UserAssetsController) QueryAssetHistory(c *gin.Context) {
	// symbol := c.Param("symbol")
	//测试不处理userid
	// userId := c.Query("userId")

	ctx := context.Background()
	assetLogs, err := ctrl.repo.FindAssetHistory(ctx)
	if err != nil {
		ctrl.logger.Error("Query", zap.Error(err))
	}
	c.JSON(http.StatusOK, gin.H{
		"data": assetLogs,
		"ok":   true,
	})
}
