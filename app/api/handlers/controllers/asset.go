package controllers

import (
	"context"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	crud_types "github.com/duolacloud/crud-core/types"
	"github.com/yzimhao/trading_engine/v2/app/common"
	"github.com/yzimhao/trading_engine/v2/internal/models/asset"
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
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
		common.ResponseError(c, err)
		return
	}

	ctx := context.Background()
	transId := uuid.New().String()
	if err := ctrl.repo.Despoit(ctx, transId, req.UserId, req.Symbol, types.Numeric(req.Amount)); err != nil {
		common.ResponseError(c, err)
		return
	}

	common.ResponseOK(c, transId)

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
		common.ResponseError(c, err)
		return
	}

	ctx := context.Background()
	transId := uuid.New().String()
	if err := ctrl.repo.Withdraw(ctx, transId, req.UserId, req.Symbol, types.Numeric(req.Amount)); err != nil {
		common.ResponseError(c, err)
		return
	}

	common.ResponseOK(c, transId)
}

// @Summary get wallet assets
// @Description get assets balance
// @ID v1.asset.query
// @Tags asset
// @Accept json
// @Produce json
// @Param symbols query string true "symbols example: BTC,ETH,USDT"
// @Success 200 {object} []asset.Asset
// @Router /api/v1/asset/query [get]
func (ctrl *UserAssetsController) Query(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userId := claims["userId"].(string)

	symbols := c.Query("symbols")
	symbolsSlice := strings.Split(symbols, ",")

	var assets []*asset.Asset

	ctx := context.Background()
	assets, err := ctrl.repo.Query(ctx, &crud_types.PageQuery{
		Filter: map[string]any{
			"symbol": map[string]any{
				"in": symbolsSlice,
			},
			"user_id": map[string]any{
				"eq": userId,
			},
		},
	})
	if err != nil {
		common.ResponseError(c, err)
		return
	}

	common.ResponseOK(c, assets)
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
		common.ResponseError(c, err)
		return
	}

	ctx := context.Background()
	transId := uuid.New().String()
	if err := ctrl.repo.Transfer(ctx, transId, req.From, req.To, req.Symbol, types.Numeric(req.Amount)); err != nil {
		common.ResponseError(c, err)
		return
	}

	common.ResponseOK(c, transId)
}

// @Summary get asset history
// @Description get an asset history
// @ID v1.asset.history
// @Tags asset
// @Accept json
// @Produce json
// @Param symbol path string true "symbol"
// @Success 200
// @Router /api/v1/asset/{symbol}/history [get]
func (ctrl *UserAssetsController) QueryAssetHistory(c *gin.Context) {
	// symbol := c.Param("symbol")
	//测试不处理userid
	// userId := c.Query("userId")

	// ctx := context.Background()
	// assetLogs, err := ctrl.repo.FindAssetHistory(ctx)
	// if err != nil {
	// 	common.ResponseError(c, err)
	// 	return
	// }
	// common.ResponseOK(c, assetLogs)
}
