package provider

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewGin(v *viper.Viper, logger *zap.Logger) *gin.Engine {
	engine := gin.New()
	return engine
}
