package di

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type HttpServer struct {
	addr   string
	engine *gin.Engine
	logger *zap.Logger
}

func NewHttpServer(in context.Context, v *viper.Viper) Server {
	v.SetDefault("port", 13081)
	port := v.GetInt("port")

	registerFn := func(s *gin.RouteInfo) {

	}
}
