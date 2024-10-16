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
	repo   persistence.AssetsRepository
	logger *zap.Logger
}

func NewUserAssetsController(repo persistence.AssetsRepository, logger *zap.Logger) *UserAssetsController {
	return &UserAssetsController{repo: repo, logger: logger}
}

type DespoitRequest struct {
	UserId string `json:"user_id"`
	Symbol string `json:"symbol"`
	Amount string `json:"amount"`
}

// @Summary asset despoit
// @Description despoit an asset
// @ID v1.wallet.asset.despoit
// @Tags wallet
// @Accept json
// @Produce json
// @Param body body DespoitRequest true "despoit request"
// @Success 200 {string} order_id
// @Router /api/v1/wallet/assets/despoit [post]
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
// @ID v1.wallet.asset.withdraw
// @Tags wallet
// @Accept json
// @Produce json
// @Param body body WithdrawRequest true "withdraw request"
// @Success 200 {string} order_id
// @Router /api/v1/wallet/assets/withdraw [post]
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
// @ID v1.wallet.asset.query
// @Tags wallet
// @Accept json
// @Produce json
// @Param symbol path string true "symbol"
// @Query userId query string true "userId测试用参数"
// @Success 200 {object} entities.Assets
// @Router /api/v1/wallet/assets/{symbol} [get]
func (ctrl *UserAssetsController) Query(c *gin.Context) {

	symbol := c.Param("symbol")
	userId := c.Query("userId")

	var asset *entities.Assets

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
// @ID v1.wallet.asset.transfer
// @Tags wallet
// @Accept json
// @Produce json
// @Param symbol path string true "symbol"
// @Param body body TransferRequest true "transfer request"
// @Router /api/v1/wallet/transfer/{symbol} [post]
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

// @Summary get wallet asset history
// @Description get an asset history
// @ID v1.wallet.asset.history
// @Tags wallet
// @Accept json
// @Produce json
// @Param symbol path string true "symbol"
// @Success 200 {object} []entities.Assets
// @Router /api/v1/wallet/assets/{symbol}/history [get]
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
