package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/yzimhao/trading_engine/v2/app/common"
	"github.com/yzimhao/trading_engine/v2/app/middlewares"
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ExampleController struct {
	engine *gin.Engine
	logger *zap.Logger
	asset  persistence.AssetRepository
	auth   *middlewares.AuthMiddleware
}

type inContext struct {
	fx.In
	Engine *gin.Engine
	Logger *zap.Logger
	Asset  persistence.AssetRepository
	Auth   *middlewares.AuthMiddleware
}

func NewExampleController(in inContext) *ExampleController {
	example := ExampleController{
		engine: in.Engine,
		logger: in.Logger,
		asset:  in.Asset,
		auth:   in.Auth,
	}

	example.registerRoutes()
	return &example
}

func (exa *ExampleController) registerRoutes() {

	exampleGroup := exa.engine.Group("example")
	exampleGroup.GET("/", exa.example)
	exampleGroup.GET("/:symbol", exa.example)
	exampleGroup.GET("/deposit", exa.auth.Auth(), exa.deposit)
}

func (exa *ExampleController) example(ctx *gin.Context) {

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

func (exa *ExampleController) deposit(ctx *gin.Context) {
	userId := common.GetUserId(ctx)

	symbols := []string{"usdt", "jpy", "eur", "btc"}

	for _, symbol := range symbols {
		transId := time.Now().Format("20060102")
		if err := exa.asset.Despoit(ctx, "deposit."+symbol+"."+transId, userId, symbol, types.Numeric("1000")); err != nil {
			common.ResponseError(ctx, err)
			exa.logger.Error("deposit error", zap.Error(err))
		}
	}

	common.ResponseOK(ctx, "success")
}
