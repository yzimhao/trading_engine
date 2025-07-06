package example

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"

	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/modules/middlewares"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type exampleModule struct {
	router    *provider.Router
	logger    *zap.Logger
	userAsset persistence.UserAssetRepository
	auth      *middlewares.AuthMiddleware
}

type inContext struct {
	fx.In
	Router    *provider.Router
	Logger    *zap.Logger
	UserAsset persistence.UserAssetRepository
	Auth      *middlewares.AuthMiddleware
}

func newExample(in inContext) {
	ex := exampleModule{
		router:    in.Router,
		logger:    in.Logger,
		userAsset: in.UserAsset,
		auth:      in.Auth,
	}
	ex.registerRoutes()
}

func (exa *exampleModule) registerRoutes() {
	exampleGroup := exa.router.Group("api/example")
	exampleGroup.GET("/", exa.example)
	exampleGroup.GET("/:symbol", exa.example)
	exampleGroup.Use(exa.auth.Auth())
	exampleGroup.GET("/deposit", exa.deposit)
}

func (exa *exampleModule) example(ctx *gin.Context) {
	support := []string{"btcusdt"}
	symbol := strings.ToLower(ctx.Param("symbol"))
	if !lo.Contains(support, symbol) {
		ctx.Redirect(301, "/example/"+support[0])
		return
	}

	ctx.HTML(http.StatusOK, "example/index.html", gin.H{
		"symbol": symbol,
	})
}

type depositReq struct {
	Asset  string          `json:"asset" binding:"required" example:"btc" form:"asset"`
	Volume decimal.Decimal `json:"volume" binding:"required" example:"btc" form:"volume"`
}

func (exa *exampleModule) deposit(ctx *gin.Context) {
	userId := exa.router.ParseUserID(ctx)

	allowSymbols := []string{"usdt", "jpy", "eur", "btc"}

	var req depositReq
	if err := ctx.BindQuery(&req); err != nil {
		exa.router.ResponseError(ctx, types.ErrInvalidParam)
		return
	}

	if !lo.Contains(allowSymbols, req.Asset) {
		exa.router.ResponseError(ctx, types.ErrInvalidParam)
		return
	}

	transId := time.Now().Format("20060102")
	if err := exa.userAsset.Despoit("auto.deposit."+req.Asset+"."+transId, userId, req.Asset, req.Volume); err != nil {
		exa.logger.Error("deposit error", zap.Error(err))
		exa.router.ResponseError(ctx, types.ErrInternalError)
		return
	}
	exa.router.ResponseOk(ctx, "")
}
