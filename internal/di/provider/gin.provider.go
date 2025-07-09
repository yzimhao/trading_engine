package provider

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewGin(v *viper.Viper, logger *zap.Logger) *gin.Engine {
	engine := gin.New()

	engine.Use(static.Serve("/", static.LocalFile("./web", false)))
	return engine
}
