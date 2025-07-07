package base

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pkgVersion "github.com/qvcloud/gopkg/version"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/yzimhao/trading_engine/v2/generated/docs"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/modules/base/asset"
	"github.com/yzimhao/trading_engine/v2/internal/modules/base/order"
	"github.com/yzimhao/trading_engine/v2/internal/modules/base/product"
	"github.com/yzimhao/trading_engine/v2/internal/modules/middlewares"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"base",
	asset.Module,
	product.Module,
	order.Module,
	fx.Invoke(run),
)

func run(router *provider.Router, authMiddleware *middlewares.AuthMiddleware) {
	registerOtherRouter(router, authMiddleware)
}

func registerOtherRouter(router *provider.Router, authMiddleware *middlewares.AuthMiddleware) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.APIv1.GET("ping", ping)
	router.APIv1.GET("version", version)

	router.APIv1.POST("/login", authMiddleware.Jwt().LoginHandler)
	router.APIv1.POST("/refresh_token", authMiddleware.Jwt().RefreshHandler)

}

// @Summary ping
// @Description test if the server is running
// @ID v1.ping
// @Tags base
// @Produce json
// @Success 200 {string} any
// @Router /api/v1/ping [get]
func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// @Summary version
// @Description 程序版本号和编译相关信息
// @ID v1.version
// @Tags base
// @Accept json
// @Produce json
// @Success 200 {string} any
// @Router /api/v1/version [get]
func version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": pkgVersion.Version,
		"go":      pkgVersion.Go,
		"build":   pkgVersion.Build,
		"commit":  pkgVersion.Commit,
	})
}
