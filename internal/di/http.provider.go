package di

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	http_server "github.com/yzimhao/trading_engine/v2/pkg/http"
)

func NewHttpServer(v *viper.Viper, engine *gin.Engine) Server {
	v.SetDefault("port", 8080)
	port := v.GetInt("port")

	return http_server.NewHttpServer(
		http_server.WithPort(port),
		http_server.WithHandler(engine),
	)

}
