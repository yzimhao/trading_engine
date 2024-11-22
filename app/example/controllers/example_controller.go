package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ExampleController struct {
	engine *gin.Engine
	logger *zap.Logger
}

type inContext struct {
	fx.In
	Engine *gin.Engine
	Logger *zap.Logger
}

func NewExampleController(in inContext) *ExampleController {
	example := ExampleController{
		engine: in.Engine,
		logger: in.Logger,
	}

	example.registerRoutes()
	return &example
}

func (exa *ExampleController) registerRoutes() {

	exampleGroup := exa.engine.Group("example")
	exampleGroup.GET("/:symbol", exa.example)
}

func (exa *ExampleController) example(ctx *gin.Context) {

	support := []string{"usdjpy", "eurusd"}
	symbol := strings.ToLower(ctx.Param("symbol"))

	if !lo.Contains(support, symbol) {
		ctx.Redirect(301, "/example/"+support[0])
		return
	}

	ctx.HTML(http.StatusOK, "example/index.html", gin.H{
		"symbol": symbol,
	})
}
