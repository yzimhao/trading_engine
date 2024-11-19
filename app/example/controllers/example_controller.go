package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ExampleController struct {
	engine *gin.Engine
}

type inContext struct {
	fx.In
	Engine *gin.Engine
}

func NewExampleController(in inContext) *ExampleController {
	example := ExampleController{
		engine: in.Engine,
	}

	example.registerRoutes()
	return &example
}

func (exa *ExampleController) registerRoutes() {
	exampleGroup := exa.engine.Group("example")
	exampleGroup.GET("/", exa.example)
}

func (exa *ExampleController) example(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "example"})
}
