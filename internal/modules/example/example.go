package example

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/v2/app/common"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/modules/middlewares"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
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
	exampleGroup := exa.router.Group("example")
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

func (exa *exampleModule) deposit(ctx *gin.Context) {
	userId := common.GetUserId(ctx)

	symbols := []string{"usdt", "jpy", "eur", "btc"}

	for _, symbol := range symbols {
		transId := time.Now().Format("20060102")
		if err := exa.userAsset.Despoit("deposit."+symbol+"."+transId, userId, symbol, decimal.NewFromFloat(1000)); err != nil {
			common.ResponseError(ctx, err)
			exa.logger.Error("deposit error", zap.Error(err))
		}
	}

	common.ResponseOK(ctx, "success")
}
