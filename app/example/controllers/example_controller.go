package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	exampleGroup.GET("/", exa.example)
}

func (exa *ExampleController) example(ctx *gin.Context) {
	exa.logger.Info("example")
	ctx.HTML(http.StatusOK, "example/index.html", gin.H{})
}
